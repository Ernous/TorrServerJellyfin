package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"server/dlna"
	"server/log"
	set "server/settings"
	"server/tmdb"
	"server/torr"
	"server/torr/state"
	"server/utils"
	apiutils "server/web/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Action: add, get, set, rem, list, drop
type torrReqJS struct {
	requestI
	Link     string `json:"link,omitempty"`
	Hash     string `json:"hash,omitempty"`
	Title    string `json:"title,omitempty"`
	Category string `json:"category,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Data     string `json:"data,omitempty"`
	StrmDir  string `json:"strm_dir,omitempty"`  // Custom directory for .strm files
	SaveToDB bool   `json:"save_to_db,omitempty"`
}

// torrents godoc
//
//	@Summary		Handle torrents informations
//	@Description	Allow to list, add, remove, get, set, drop, wipe torrents on server. The action depends of what has been asked.
//
//	@Tags			API
//
//	@Param			request	body	torrReqJS	true	"Torrent request. Available params for action: add, get, set, rem, list, drop, wipe. link required for add, hash required for get, set, rem, drop."
//
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/torrents [post]
func torrents(c *gin.Context) {
	var req torrReqJS
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusBadRequest)
	switch req.Action {
	case "add":
		{
			addTorrent(req, c)
		}
	case "get":
		{
			getTorrent(req, c)
		}
	case "set":
		{
			setTorrent(req, c)
		}
	case "rem":
		{
			remTorrent(req, c)
		}
	case "list":
		{
			listTorrents(c)
		}
	case "drop":
		{
			dropTorrent(req, c)
		}
	case "wipe":
		{
			wipeTorrents(c)
		}
	}
}

func addTorrent(req torrReqJS, c *gin.Context) {
	if req.Link == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("link is empty"))
		return
	}

	log.TLogln("add torrent", req.Link)
	req.Link = strings.ReplaceAll(req.Link, "&amp;", "&")
	torrSpec, err := apiutils.ParseLink(req.Link)
	if err != nil {
		log.TLogln("error parse link:", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Try to fetch metadata from TMDB if API key is configured
	if set.BTsets.TMDBApiKey != "" && req.Title == "" {
		tmdbClient := tmdb.NewClient(set.BTsets.TMDBApiKey)
		media, err := tmdbClient.SearchAuto(torrSpec.DisplayName)
		if err == nil {
			req.Title = media.GetTitle()
			req.Poster = media.GetPosterURL()
			log.TLogln("TMDB: Fetched metadata for:", req.Title)
		}
	}

	tor, err := torr.AddTorrent(torrSpec, req.Title, req.Poster, req.Data, req.Category, req.StrmDir)
	// if tor.Data != "" && set.BTsets.EnableDebug {
	// 	log.TLogln("torrent data:", tor.Data)
	// }
	// if tor.Category != "" && set.BTsets.EnableDebug {
	// 	log.TLogln("torrent category:", tor.Category)
	// }
	if err != nil {
		log.TLogln("error add torrent:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		if !tor.GotInfo() {
			log.TLogln("error add torrent:", "timeout connection get torrent info")
			return
		}

		if tor.Title == "" {
			tor.Title = torrSpec.DisplayName // prefer dn over name
			tor.Title = strings.ReplaceAll(tor.Title, "rutor.info", "")
			tor.Title = strings.ReplaceAll(tor.Title, "_", " ")
			tor.Title = strings.Trim(tor.Title, " ")
			if tor.Title == "" {
				tor.Title = tor.Name()
			}
		}

		if req.SaveToDB {
			torr.SaveTorrentToDB(tor)
		}

		// Auto-create .strm files if enabled
		if set.BTsets.JlfnAutoCreate && set.BTsets.JlfnAddr != "" {
			log.TLogln("Auto-creating .strm files for torrent:", tor.Hash().HexString())
			createStrmFilesForTorrent(tor, c)
		}
	}()
	// TODO: remove
	if set.BTsets.EnableDLNA {
		dlna.Stop()
		dlna.Start()
	}
	c.JSON(200, tor.Status())
}

func getTorrent(req torrReqJS, c *gin.Context) {
	if req.Hash == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("hash is empty"))
		return
	}
	tor := torr.GetTorrent(req.Hash)

	if tor != nil {
		st := tor.Status()
		c.JSON(200, st)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func setTorrent(req torrReqJS, c *gin.Context) {
	if req.Hash == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("hash is empty"))
		return
	}
	torr.SetTorrent(req.Hash, req.Title, req.Poster, req.Category, req.Data)
	c.Status(200)
}

func remTorrent(req torrReqJS, c *gin.Context) {
	if req.Hash == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("hash is empty"))
		return
	}
	
	// Get torrent info before removing for .strm cleanup
	tor := torr.GetTorrent(req.Hash)
	if tor != nil && set.BTsets.JlfnAddr != "" {
		removeStrmFilesForTorrent(tor)
	}
	
	torr.RemTorrent(req.Hash)
	// TODO: remove
	if set.BTsets.EnableDLNA {
		dlna.Stop()
		dlna.Start()
	}
	c.Status(200)
}

func listTorrents(c *gin.Context) {
	list := torr.ListTorrent()
	if len(list) == 0 {
		c.JSON(200, []*state.TorrentStatus{})
		return
	}
	var stats []*state.TorrentStatus
	for _, tr := range list {
		stats = append(stats, tr.Status())
	}
	c.JSON(200, stats)
}

func dropTorrent(req torrReqJS, c *gin.Context) {
	if req.Hash == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("hash is empty"))
		return
	}
	torr.DropTorrent(req.Hash)
	c.Status(200)
}

func wipeTorrents(c *gin.Context) {
	torrents := torr.ListTorrent()
	for _, t := range torrents {
		torr.RemTorrent(t.TorrentSpec.InfoHash.HexString())
	}
	// TODO: remove (copied todo from remTorrent())
	if set.BTsets.EnableDLNA {
		dlna.Stop()
		dlna.Start()
	}
	c.Status(200)
}

// torrentsJlfn godoc
//
//	@Summary		Add torrent for Jellyfin with .strm files creation
//	@Description	Add torrent and automatically create .strm files for Jellyfin
//	@Tags			API
//	@Param			request	body	torrReqJS	true	"Torrent request with link"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/torrents/jlfn [post]
func torrentsJlfn(c *gin.Context) {
	var req torrReqJS
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	addJlfn(req, c)
}

// createStrmFilesForTorrent creates .strm files for a torrent
func createStrmFilesForTorrent(tor *torr.Torrent, c *gin.Context) {
	// Wait a bit for torrent to initialize
	time.Sleep(5 * time.Second)

	// Check base path configuration
	basePath := set.BTsets.JlfnAddr
	if basePath == "" {
		log.TLogln("JlfnAddr not configured, skipping .strm creation")
		return
	}

	// Load torrent from DB if needed
	if tor.Stat == state.TorrentInDB {
		tor = torr.LoadTorrent(tor)
		if tor == nil {
			log.TLogln("Failed to load torrent from DB")
			return
		}
	}

	// Check if torrent's underlying Torrent object is valid
	if tor.Torrent == nil {
		log.TLogln("Torrent object is nil, cannot create .strm files")
		return
	}

	// Get host for stream URLs
	// Use configured TorrServerHost if available, otherwise use request host
	host := set.BTsets.TorrServerHost
	if host == "" {
		host = utils.GetScheme(c) + "://" + c.Request.Host
		log.TLogln("TorrServerHost not configured, using request host:", host)
	} else {
		log.TLogln("Using configured TorrServerHost:", host)
	}
	torrents := tor.Status()

	// Determine category (serials/films)
	catPath := ""
	torName := tor.Name()
	
	// Use torrent category from web UI if available
	torCategory := strings.ToLower(tor.Category)
	log.TLogln("Torrent category from UI:", torCategory)
	
	// Map Russian/English categories to folder names
	switch torCategory {
	case "сериалы", "серіали", "serials", "series", "tv", "tv shows":
		catPath = "torrSerials"
		log.TLogln("Category mapped to torrSerials (from UI category)")
	case "фильмы", "фільми", "movies", "films":
		catPath = "torrFilms"
		log.TLogln("Category mapped to torrFilms (from UI category)")
	default:
		// Auto-detect if category not set or unknown
		torTitle := strings.ToLower(torName)
		seasonPattern := regexp.MustCompile(`(?i)s\d+|season\s+\d+`)
		episodePattern := regexp.MustCompile(`(?i)e\d+|episode\s+\d+`)
		
		isSeries := seasonPattern.MatchString(torTitle) || episodePattern.MatchString(torTitle)
		
		if isSeries {
			catPath = "torrSerials"
			log.TLogln("Auto-detected as TV series (found season/episode pattern)")
		} else {
			catPath = "torrFilms"
			log.TLogln("Auto-detected as movie (no season/episode pattern)")
		}
	}

	log.TLogln("Final category path:", catPath, "for:", torName)

	// Create full path: basePath/category/custom_dir/torrent_name
	// If custom directory is specified, use it between category and torrent name
	var fullBasePath string
	if tor.StrmDir != "" {
		fullBasePath = filepath.Join(basePath, catPath, tor.StrmDir, torName)
		log.TLogln("Using custom directory:", tor.StrmDir)
	} else {
		fullBasePath = filepath.Join(basePath, catPath, torName)
	}
	log.TLogln("Creating .strm files in:", fullBasePath)

	// Create .strm files for each video file
	for _, f := range torrents.FileStats {
		if utils.GetMimeType(f.Path) != "*/*" {
			// Create stream URL
			fileName := filepath.Base(f.Path)
			streamURL := fmt.Sprintf("%s/stream/%s?link=%s&index=%d&play",
				host,
				url.PathEscape(fileName),
				torrents.Hash,
				f.Id)

			// Create .strm filename
			strmName := filepath.Base(f.Path)
			strmName = strings.ReplaceAll(strmName, "/", "")
			strmName = strings.ReplaceAll(strmName, "\\", "")

			// Remove original extension and add .strm
			ext := filepath.Ext(strmName)
			if ext != "" {
				strmName = strings.TrimSuffix(strmName, ext)
			}
			strmName += ".strm"

			// Full path to .strm file
			strmPath := filepath.Join(fullBasePath, strmName)

			// Create directory
			if err := os.MkdirAll(filepath.Dir(strmPath), 0755); err != nil {
				log.TLogln("Error creating directory:", err)
				continue
			}

			// Write .strm file
			if err := os.WriteFile(strmPath, []byte(streamURL), 0644); err != nil {
				log.TLogln("Error writing .strm file:", err)
				continue
			}

			log.TLogln("Created .strm file:", strmPath)
		}
	}

	log.TLogln("Finished creating .strm files for:", torrents.Hash)

	// Save the path to database
	tor.StrmPath = fullBasePath
	torr.SaveTorrentToDB(tor)
	log.TLogln("Saved .strm path to DB:", fullBasePath)
}

// removeStrmFilesForTorrent removes .strm files and directories for a torrent
func removeStrmFilesForTorrent(tor *torr.Torrent) {
	// Use saved path if available
	if tor.StrmPath != "" {
		log.TLogln("Removing .strm files from saved path:", tor.StrmPath)
		
		// Remove torrent directory
		if err := os.RemoveAll(tor.StrmPath); err != nil {
			log.TLogln("Error removing torrent directory:", err)
		} else {
			log.TLogln("Removed torrent directory:", tor.StrmPath)
		}

		// Try to remove parent custom directory if empty
		parentDir := filepath.Dir(tor.StrmPath)
		basePath := set.BTsets.JlfnAddr
		if basePath != "" {
			// Don't remove if it's a category folder (torrSerials/torrFilms)
			if !strings.HasSuffix(parentDir, "torrSerials") && !strings.HasSuffix(parentDir, "torrFilms") {
				entries, err := os.ReadDir(parentDir)
				if err == nil && len(entries) == 0 {
					if err := os.Remove(parentDir); err != nil {
						log.TLogln("Error removing parent directory:", err)
					} else {
						log.TLogln("Removed empty parent directory:", parentDir)
					}
				}
			}
		}
		return
	}

	// Fallback: calculate path from torrent data
	basePath := set.BTsets.JlfnAddr
	if basePath == "" {
		return
	}

	torName := tor.Name()
	log.TLogln("Removing .strm files for:", torName, "(no saved path, calculating)")

	// Determine category (same logic as in createStrmFilesForTorrent)
	catPath := ""
	torCategory := strings.ToLower(tor.Category)
	
	switch torCategory {
	case "сериалы", "серіали", "serials", "series", "tv", "tv shows":
		catPath = "torrSerials"
	case "фильмы", "фільми", "movies", "films":
		catPath = "torrFilms"
	default:
		// Auto-detect
		torTitle := strings.ToLower(torName)
		seasonPattern := regexp.MustCompile(`(?i)s\d+|season\s+\d+`)
		episodePattern := regexp.MustCompile(`(?i)e\d+|episode\s+\d+`)
		
		isSeries := seasonPattern.MatchString(torTitle) || episodePattern.MatchString(torTitle)
		
		if isSeries {
			catPath = "torrSerials"
		} else {
			catPath = "torrFilms"
		}
	}

	// Build path to torrent directory
	var torrentDirPath string
	if tor.StrmDir != "" {
		torrentDirPath = filepath.Join(basePath, catPath, tor.StrmDir, torName)
	} else {
		torrentDirPath = filepath.Join(basePath, catPath, torName)
	}

	// Remove torrent directory
	if err := os.RemoveAll(torrentDirPath); err != nil {
		log.TLogln("Error removing torrent directory:", err)
	} else {
		log.TLogln("Removed torrent directory:", torrentDirPath)
	}

	// If custom directory was used, try to remove it if empty
	if tor.StrmDir != "" {
		customDirPath := filepath.Join(basePath, catPath, tor.StrmDir)
		entries, err := os.ReadDir(customDirPath)
		if err == nil && len(entries) == 0 {
			if err := os.Remove(customDirPath); err != nil {
				log.TLogln("Error removing custom directory:", err)
			} else {
				log.TLogln("Removed empty custom directory:", customDirPath)
			}
		}
	}
}

func addJlfn(req torrReqJS, c *gin.Context) {
	// First, add torrent normally
	addTorrent(req, c)

	// Then create .strm files asynchronously
	go func() {
		// Wait for torrent to initialize
		time.Sleep(15 * time.Second)

		// Parse link to get hash
		req.Link = strings.ReplaceAll(req.Link, "&amp;", "&")
		torrSpec, err := apiutils.ParseLink(req.Link)
		if err != nil {
			log.TLogln("error parse link:", err)
			return
		}

		hash := torrSpec.InfoHash.String()
		log.TLogln("Creating .strm files for torrent:", hash)

		// Get torrent
		tor := torr.GetTorrent(hash)
		if tor == nil {
			log.TLogln("Torrent not found:", hash)
			return
		}

		// Create .strm files
		createStrmFilesForTorrent(tor, c)

		// Optionally drop torrent after delay
		go func() {
			time.Sleep(15 * time.Second)
			torr.DropTorrent(hash)
			log.TLogln("Dropped torrent:", hash)
		}()
	}()
}
