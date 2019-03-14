using System;
using System.Threading.Tasks;
using StackExchange.Redis;

namespace menunotifier
{
    class Program
    {
        static async Task Main(string[] args)
        {
            System.Console.WriteLine("Starting up...");

            int waittime = int.Parse(Environment.GetEnvironmentVariable("SLEEP_INTERVAL"));

            var redis = await SetupRedis();
            var checker = new Checker(redis);

            System.Console.WriteLine("Starting monitoring!");

            while (true)
            {
                System.Console.WriteLine("Checking menu changed...");
                if (await checker.CheckMenuChanged())
                    System.Console.WriteLine("Menu changed!");
                await Task.Delay(1000 * waittime);
            }
        }

        static async Task<ConnectionMultiplexer> SetupRedis()
        {
            return await ConnectionMultiplexer.ConnectAsync(new ConfigurationOptions()
            {
                EndPoints = 
                {
                    { Environment.GetEnvironmentVariable("REDIS_HOST"), 6379 }
                },
                Password = Environment.GetEnvironmentVariable("REDIS_PASS"),
            });
        }
    }
}
