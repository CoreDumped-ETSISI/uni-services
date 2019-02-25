using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Google.Apis.Services;
using Google.Apis.Sheets.v4;
using System.Linq;
using Newtonsoft.Json;
using bibliosalas.Model;
using Microsoft.Extensions.Hosting;
using System.Threading;
using Microsoft.Extensions.Logging;

namespace bibliosalas.Services
{
    public class CacheService : BackgroundService
    {
        private static Biblioteca _biblioteca;
        public static Biblioteca Biblioteca
        {
            get
            {
                return _biblioteca;
            }
            private set
            {
                _biblioteca = value;
            }
        }

        private readonly ILogger<CacheService> _log;
        private readonly int UpdateTime;
        private readonly SalasService _salas;

        public CacheService(ILogger<CacheService> log, SalasService salas)
        {
            _log = log;
            _salas = salas;

            UpdateTime = 15 * 60; // 15 minutes

            string updatevar = Environment.GetEnvironmentVariable("CACHE_UPDATE_INTERVAL");
            if (!string.IsNullOrEmpty(updatevar) && int.TryParse(updatevar, out int newinterval))
            {
                UpdateTime = newinterval;
            }
        }

        private async Task FetchData()
        {
            _log.LogInformation("Actualizando cache...");

            try
            {
                Biblioteca = await _salas.GetSalaState();
            }
            catch (Exception ex)
            {
                _log.LogError(ex, "Error al actualizar cache!");
            }
        }

        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            _log.LogDebug("CacheService is starting");

            stoppingToken.Register(() => _log.LogDebug("CacheService is stopping"));

            while (!stoppingToken.IsCancellationRequested)
            {
                await FetchData();
                await Task.Delay(UpdateTime * 1000, stoppingToken);
            }
        }
    }
}
