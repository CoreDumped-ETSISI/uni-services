using System.Collections.Generic;
using github.Model;
using System.Net.Http;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace github.Api
{
    public static class GitHub
    {
        public static async Task<List<Repository>> GetOrganizationRepositories(string org)
        {
            string url = string.Format("https://api.github.com/orgs/{0}/repos", org);

            using (var wc = new HttpClient())
            {
                wc.DefaultRequestHeaders.Add("User-Agent", "guad");
                string raw = await wc.GetStringAsync(url);
                var repos = JsonConvert.DeserializeObject<List<Repository>>(raw);
                return repos;
            }
        }

        public static async Task<Dictionary<string, int>> GetRepositoryLanguages(string repo)
        {
            string url = string.Format("https://api.github.com/repos/{0}/languages", repo);

            using (var wc = new HttpClient())
            {
                wc.DefaultRequestHeaders.Add("User-Agent", "guad");
                string raw = await wc.GetStringAsync(url);
                var repos = JsonConvert.DeserializeObject<Dictionary<string, int>>(raw);
                return repos;
            }
        }
    }
}