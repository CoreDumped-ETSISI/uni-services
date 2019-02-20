using System.Collections.Generic;
using Newtonsoft.Json;

namespace bibliosalas.Model
{
    public class SalaTimesheet
    {
        public int Id { get; set; }
        public List<TimeSlot> Occupied { get; set; }

        [JsonIgnore]
        public bool[] OccupiedMap { get; set; }

        public void GenerateTimeslots()
        {
            Occupied = new List<TimeSlot>();

            int? start = null;
            for (int i = 0; i < OccupiedMap.Length; i++)
            {
                bool occupied = OccupiedMap[i];

                if (occupied && start == null)
                {
                    start = i;
                }
                else if (!occupied && start != null)
                {
                    Occupied.Add(new TimeSlot()
                    {
                        Start = TimeSlot.IndexToTime(start.Value),
                        End = TimeSlot.IndexToTime(i),
                    });

                    start = null;
                }
            }

            if (start != null)
            {
                Occupied.Add(new TimeSlot()
                {
                    Start = TimeSlot.IndexToTime(start.Value),
                    End = TimeSlot.IndexToTime(OccupiedMap.Length),
                });
            }
        }
    }
}