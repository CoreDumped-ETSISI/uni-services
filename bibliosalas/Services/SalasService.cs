using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Google.Apis.Services;
using Google.Apis.Sheets.v4;
using System.Linq;
using Newtonsoft.Json;
using bibliosalas.Model;

namespace bibliosalas.Services
{
    public class SalasService
    {
        private readonly GoogleSheetsService _sheets;
        private readonly string SheetId;
        private readonly string SheetRange;

        public SalasService(GoogleSheetsService sheets)
        {
            _sheets = sheets;
            SheetId = Environment.GetEnvironmentVariable("SHEET_ID");
            SheetRange = Environment.GetEnvironmentVariable("SHEET_RANGE");
        }

        private async Task<Biblioteca> FetchData()
        {
            var sheet = await _sheets.GetSheet(SheetId, SheetRange);

            var biblio = new Biblioteca();

            biblio.Salas = new List<SalaTimesheet>
            {
                new SalaTimesheet() { Id = 1, OccupiedMap = new bool[22], },
                new SalaTimesheet() { Id = 2, OccupiedMap = new bool[22], },
                new SalaTimesheet() { Id = 3, OccupiedMap = new bool[22], },
                new SalaTimesheet() { Id = 4, OccupiedMap = new bool[22], },
                new SalaTimesheet() { Id = 5, OccupiedMap = new bool[22], },
            };

            var rows = sheet.Sheets.First().Data.First().RowData;
            for (var i = 0; i < rows.Count; i++)
            {
                var row = rows[i];

                // Columns
                for (int j = 0; j < row.Values.Count; j++)
                {
                    var cell = row.Values[j];

                    var r = cell.EffectiveFormat.BackgroundColor.Red;
                    var g = cell.EffectiveFormat.BackgroundColor.Green;
                    var b = cell.EffectiveFormat.BackgroundColor.Blue;

                    bool free = r == g && g == b;
                    biblio.Salas[j].OccupiedMap[i] = !free;
                }
            }

            foreach (var sala in biblio.Salas)
                sala.GenerateTimeslots();

            return biblio;
        }

        public async Task<Biblioteca> GetSalaState()
        {
            return await FetchData();
        }
        
    }
}