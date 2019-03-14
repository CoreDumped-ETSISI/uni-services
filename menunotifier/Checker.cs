using System;
using StackExchange.Redis;
using Newtonsoft.Json;
using menunotifier.Model;
using System.Net.Http;
using System.Threading.Tasks;

namespace menunotifier
{
    public class Checker
    {
        public Checker(ConnectionMultiplexer redis)
        {
            _redis = redis;
        }

        private readonly ConnectionMultiplexer _redis;
        private const string LastFetchKey = "CAFETERIA_MENU_LAST_FETCH_URL";
        private const string NotificationChannel = "CAFETERIA_NEW_MENU_AVAILABLE";

        public async Task<bool> CheckMenuChanged()
        {
            try
            {                
                CafetApiResponse resp = await Api.FetchMenu();

                IDatabase db = _redis.GetDatabase(0, new object());
                string lasturl = await db.StringGetAsync(LastFetchKey);
                await db.StringSetAsync(LastFetchKey, resp.FetchURL);

                if (lasturl == null || lasturl == resp.FetchURL)
                    return false;

                await Notify(resp);
                
                return true;
            }
            catch (HttpRequestException ex)
            {
                System.Console.WriteLine("Could not fetch menu:");
                System.Console.WriteLine(ex);
                return false;
            }
        }

        private async Task Notify(CafetApiResponse resp)
        {
            ISubscriber sub = _redis.GetSubscriber();

            RedisNotification notif = new RedisNotification()
            {
                Text = "El menú de esta semana ya está disponible! Utiliza el comando /menu para mirar el de hoy.",
                Link = resp.FetchURL,
            };

            await sub.PublishAsync(NotificationChannel, JsonConvert.SerializeObject(notif));
        }
    }
}