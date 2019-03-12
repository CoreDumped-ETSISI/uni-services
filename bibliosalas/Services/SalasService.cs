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

                    float r, g, b;

                    if (cell.EffectiveFormat != null && cell.EffectiveFormat.BackgroundColor != null)
                    {
                        r = cell.EffectiveFormat.BackgroundColor.Red.Value;
                        g = cell.EffectiveFormat.BackgroundColor.Green.Value;
                        b = cell.EffectiveFormat.BackgroundColor.Blue.Value;
                    }
                    else if (cell.UserEnteredFormat != null && cell.UserEnteredFormat.BackgroundColor != null)
                    {
                        r = cell.UserEnteredFormat.BackgroundColor.Red.Value;
                        g = cell.UserEnteredFormat.BackgroundColor.Green.Value;
                        b = cell.UserEnteredFormat.BackgroundColor.Blue.Value;
                    }
                    else
                    {
                        r = 1;
                        g = 1;
                        b = 1;
                    }


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