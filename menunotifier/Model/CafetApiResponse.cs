using System;

namespace menunotifier.Model
{
    public class CafetApiResponse
    {
        public FetchMetadata LastFetch { get; set; }
        public DateTime LastFetched { get; set; }
        public string FetchURL { get; set; }
    }
}