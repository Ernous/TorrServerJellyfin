using MediaBrowser.Model.Plugins;

namespace Jellyfin.Plugin.TorrServer.Configuration
{
    public class PluginConfiguration : BasePluginConfiguration
    {
        public string TorrServerUrl { get; set; } = "http://localhost:8090";
        
        public string TorrServerUsername { get; set; } = "";
        
        public string TorrServerPassword { get; set; } = "";
        
        public bool AutoCreateStrm { get; set; } = true;
        
        public string StrmPath { get; set; } = "/media/strm";
        
        public string DefaultCategory { get; set; } = "Movies";
        
        public int CacheSize { get; set; } = 64;
        
        public int ReaderReadAHead { get; set; } = 95;
        
        public int PreloadCache { get; set; } = 50;
        
        public bool UseDisk { get; set; } = true;
        
        public string DiskCachePath { get; set; } = "/cache/torrserver";
        
        public bool RemoveCacheOnDrop { get; set; } = false;
        
        public bool RemoveCacheOnStop { get; set; } = false;
        
        public string TmdbApiKey { get; set; } = "";
    }
}
