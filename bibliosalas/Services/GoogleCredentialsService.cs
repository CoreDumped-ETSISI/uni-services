using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Google.Apis.Auth.OAuth2;
using Google.Apis.Http;
using Google.Apis.Sheets.v4;
using System;

namespace bibliosalas.Services
{
    public class GoogleCredentialsService : IConfigurableHttpClientInitializer, IHttpExecuteInterceptor
    {
        private string ApiKey;

        public GoogleCredentialsService()
        {
            ApiKey = System.Environment.GetEnvironmentVariable("GOOGLE_API_KEY");
        }

        public void Initialize(ConfigurableHttpClient httpClient)
        {
            httpClient.MessageHandler.AddExecuteInterceptor(this);
        }

        public Task InterceptAsync(HttpRequestMessage request, CancellationToken cancellationToken)
        {
            // The new request URI is the old one plus a new query parameter (key)
            UriBuilder builder = new UriBuilder(request.RequestUri);
            string oldq = builder.Query;

            if (string.IsNullOrEmpty(oldq))
            {
                oldq = "key=" + this.ApiKey;
            }
            else
            {
                oldq = oldq.TrimStart('?') + "&key=" + this.ApiKey;
            }

            builder.Query = oldq;

            request.RequestUri = builder.Uri;

            return Task.FromResult(0);
        }
    }
}