namespace Jellyfin.Plugin.TorrServer.Models
{
    public class AddTorrentRequest
    {
        public string Link { get; set; } = string.Empty;
        public string Title { get; set; } = string.Empty;
        public string Poster { get; set; } = string.Empty;
        public string Data { get; set; } = string.Empty;
        public string Category { get; set; } = "Movies";
        public string StrmDir { get; set; } = string.Empty;
        public bool CreateStrm { get; set; } = true;
    }

    public class RemoveTorrentRequest
    {
        public string Hash { get; set; } = string.Empty;
    }

    public class UpdateSettingsRequest
    {
        public int CacheSize { get; set; }
        public int ReaderReadAHead { get; set; }
        public int PreloadCache { get; set; }
        public bool UseDisk { get; set; }
        public string TorrentsSavePath { get; set; } = string.Empty;
        public bool RemoveCacheOnDrop { get; set; }
        public string JlfnAddr { get; set; } = string.Empty;
        public bool JlfnAutoCreate { get; set; }
        public string TorrServerHost { get; set; } = string.Empty;
        public string TMDBApiKey { get; set; } = string.Empty;
    }

    public class TmdbSearchRequest
    {
        public string Query { get; set; } = string.Empty;
        public string Language { get; set; } = "en";
        public string Type { get; set; } = "multi";
    }
}
