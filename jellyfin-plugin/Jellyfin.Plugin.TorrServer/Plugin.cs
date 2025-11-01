using System;
using System.Collections.Generic;
using Jellyfin.Plugin.TorrServer.Configuration;
using MediaBrowser.Common.Configuration;
using MediaBrowser.Common.Plugins;
using MediaBrowser.Model.Plugins;
using MediaBrowser.Model.Serialization;

namespace Jellyfin.Plugin.TorrServer
{
    public class Plugin : BasePlugin<PluginConfiguration>, IHasWebPages
    {
        public override string Name => "TorrServer";

        public override Guid Id => Guid.Parse("a8b12c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d");

        public Plugin(IApplicationPaths applicationPaths, IXmlSerializer xmlSerializer)
            : base(applicationPaths, xmlSerializer)
        {
            Instance = this;
        }

        public static Plugin? Instance { get; private set; }

        public IEnumerable<PluginPageInfo> GetPages()
        {
            return new[]
            {
                new PluginPageInfo
                {
                    Name = this.Name,
                    EmbeddedResourcePath = string.Format("{0}.Configuration.configPage.html", GetType().Namespace)
                },
                new PluginPageInfo
                {
                    Name = "torrentsPage",
                    EmbeddedResourcePath = string.Format("{0}.Configuration.torrentsPage.html", GetType().Namespace)
                }
            };
        }
    }
}
