using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Jellyfin.Plugin.TorrServer.Models;
using MediaBrowser.Controller;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace Jellyfin.Plugin.TorrServer.Api
{
    [ApiController]
    [Authorize(Policy = "RequiresElevation")]
    [Route("TorrServer")]
    public class TorrServerController : ControllerBase
    {
        private readonly ILogger<TorrServerController> _logger;
        private readonly IHttpClientFactory _httpClientFactory;

        public TorrServerController(
            ILogger<TorrServerController> logger,
            IHttpClientFactory httpClientFactory)
        {
            _logger = logger;
            _httpClientFactory = httpClientFactory;
        }

        private string TorrServerUrl => Plugin.Instance?.Configuration.TorrServerUrl ?? "http://localhost:8090";

        [HttpGet("torrents")]
        public async Task<ActionResult> GetTorrents()
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                var response = await client.PostAsync($"{TorrServerUrl}/torrents", null);
                var content = await response.Content.ReadAsStringAsync();
                return Content(content, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error getting torrents");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("add")]
        public async Task<ActionResult> AddTorrent([FromBody] AddTorrentRequest request)
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                
                var torrentData = new
                {
                    action = "add",
                    link = request.Link,
                    title = request.Title,
                    poster = request.Poster,
                    data = request.Data,
                    save_to_db = true,
                    category = request.Category,
                    strm_dir = request.StrmDir
                };

                var json = JsonSerializer.Serialize(torrentData);
                var content = new StringContent(json, Encoding.UTF8, "application/json");
                
                var response = await client.PostAsync($"{TorrServerUrl}/torrents", content);
                var responseContent = await response.Content.ReadAsStringAsync();
                
                return Content(responseContent, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error adding torrent");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("remove")]
        public async Task<ActionResult> RemoveTorrent([FromBody] RemoveTorrentRequest request)
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                
                var torrentData = new
                {
                    action = "rem",
                    hash = request.Hash
                };

                var json = JsonSerializer.Serialize(torrentData);
                var content = new StringContent(json, Encoding.UTF8, "application/json");
                
                var response = await client.PostAsync($"{TorrServerUrl}/torrents", content);
                var responseContent = await response.Content.ReadAsStringAsync();
                
                return Content(responseContent, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error removing torrent");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpGet("settings")]
        public async Task<ActionResult> GetSettings()
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                var response = await client.GetAsync($"{TorrServerUrl}/settings");
                var content = await response.Content.ReadAsStringAsync();
                return Content(content, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error getting settings");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("settings")]
        public async Task<ActionResult> UpdateSettings([FromBody] UpdateSettingsRequest request)
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                var json = JsonSerializer.Serialize(request);
                var content = new StringContent(json, Encoding.UTF8, "application/json");
                
                var response = await client.PostAsync($"{TorrServerUrl}/settings", content);
                var responseContent = await response.Content.ReadAsStringAsync();
                
                return Content(responseContent, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error updating settings");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("search")]
        public async Task<ActionResult> SearchTmdb([FromBody] TmdbSearchRequest request)
        {
            try
            {
                var client = _httpClientFactory.CreateClient();
                var json = JsonSerializer.Serialize(request);
                var content = new StringContent(json, Encoding.UTF8, "application/json");
                
                var response = await client.PostAsync($"{TorrServerUrl}/tmdb/search", content);
                var responseContent = await response.Content.ReadAsStringAsync();
                
                return Content(responseContent, "application/json");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error searching TMDB");
                return StatusCode(500, new { error = ex.Message });
            }
        }
    }
}
