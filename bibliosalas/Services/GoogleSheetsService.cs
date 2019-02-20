using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Google.Apis.Services;
using Google.Apis.Sheets.v4;
using Newtonsoft.Json;

namespace bibliosalas.Services
{
    public class GoogleSheetsService
    {
        private readonly GoogleCredentialsService _credentials;

        public GoogleSheetsService(GoogleCredentialsService credentials)
        {
            _credentials = credentials;
        }

        public async Task<Google.Apis.Sheets.v4.Data.Spreadsheet> GetSheet(string spreadsheetId, string range)
        {
            SheetsService sheetsService = new SheetsService(new BaseClientService.Initializer
            {
                HttpClientInitializer = _credentials,
                ApplicationName = "BiblioSalas/0.1",
            });

            // The ranges to retrieve from the spreadsheet.
            List<string> ranges = new List<string> { range };
            bool includeGridData = true;

            SpreadsheetsResource.GetRequest request = sheetsService.Spreadsheets.Get(spreadsheetId);
            request.Ranges = ranges;
            request.IncludeGridData = includeGridData;

            Google.Apis.Sheets.v4.Data.Spreadsheet response = await request.ExecuteAsync();

            return response;
        }
    }
}