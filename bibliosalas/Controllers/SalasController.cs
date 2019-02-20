using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using bibliosalas.Model;
using bibliosalas.Services;
using Microsoft.AspNetCore.Mvc;

namespace bibliosalas.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    [Produces("application/json")]
    public class SalasController : ControllerBase
    {
        private readonly SalasService _salas;

        public SalasController(SalasService salas)
        {
            _salas = salas;
            
        }

        // GET api/salas
        [HttpGet]
        public object GetAllSalas()
        {
            var bib = CacheService.Biblioteca;

            return bib;
        }

        // GET api/salas/5
        [HttpGet("{id}")]
        public object GetSala(int id)
        {
            if (id < 1 || id > 5) return BadRequest();

            var bib = CacheService.Biblioteca;

            return bib.Salas[id-1];
        }

        // GET api/salas/reserve
        [HttpGet("reserve")]
        public object TryReserve(int min)
        {
            if (min <= 0) return BadRequest(new { message = "Minutos deben ser positivos" });

            var timenow = DateTime.UtcNow.AddHours(1);
            bool tomorrow = false;
            
            if (timenow.Hour >= 20)
            {
                //return BadRequest(new { message = "La biblioteca está cerrada." });
                tomorrow = true;
                timenow = new DateTime(timenow.Year, timenow.Month, timenow.Day, TimeSlot.HourStart, 0, 0);
            }

            int startindex = (timenow.Hour - TimeSlot.HourStart) * 2 + (timenow.Minute >= 30 ? 1 : 0);
            
            int steps = (int)Math.Ceiling(min / 30f);
            var bib = CacheService.Biblioteca;

            int minstart = int.MaxValue;
            int minsala = -1;
            var candidates = new List<object>();

            for (int i = 0; i < bib.Salas.Count; i++)
            {
                var sala = bib.Salas[i];
                int needed = steps;

                for (int j = startindex; j < sala.OccupiedMap.Length; j++)
                {
                    if (!sala.OccupiedMap[j])
                        needed--;
                    else needed = steps;

                    if (needed == 0)
                    {
                        int start = j - (steps - 1);
                        candidates.Add(new {
                            sala = i + 1,
                            start = TimeSlot.IndexToTime(start),
                            end = TimeSlot.IndexToTime(start + steps),
                        });

                        // Queremos la sala con el inicio menor.
                        if (start < minstart)
                        {
                            minstart = start;
                            minsala = i;
                        }
                        break;
                    }
                }
            }

            if (minsala == -1)
                return NotFound(new { found = false });

            return Ok(new {
                found = true,
                sala = minsala + 1,
                start = TimeSlot.IndexToTime(minstart),
                end = TimeSlot.IndexToTime(minstart + steps),
                tomorrow = tomorrow,
                options = candidates,
            });

        }
    }
}
