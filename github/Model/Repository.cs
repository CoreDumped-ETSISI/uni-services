using Newtonsoft.Json;

namespace github.Model
{
    public class Repository
    {
        public int Id { get; set; }
        public string Name { get; set; }

        [JsonProperty("full_name")]
        public string FullName { get; set; }

        [JsonProperty("html_url")]
        public string HtmlUrl { get; set; }
        public string Description { get; set; }
    }
}