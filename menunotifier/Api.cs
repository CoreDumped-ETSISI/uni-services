using System;
using System.Threading;
using StackExchange.Redis;
using Newtonsoft.Json;
using System.Threading.Tasks;
using menunotifier.Model;
using System.Net.Http;

namespace menunotifier
{
    public static class Api
    {
        private static HttpClient _client = new HttpClient();
        public static async Task<CafetApiResponse> FetchMenu()
        {
            var resp = await _client.GetAsync("https://cafe.kolhos.chichasov.es/");
            resp.EnsureSuccessStatusCode();

            CafetApiResponse caf = JsonConvert.DeserializeObject<CafetApiResponse>(await resp.Content.ReadAsStringAsync());

            return caf;
        }
    }
}