var site = fermata.json('');

function getCurrentStatus() {
    return new Promise(function(resolve, reject) {
        site.api.status.get(function(err, data) {
            if (err) {
                reject(err);
                return;
            }

            let pages = {};
            for (let site of data) {
                if (site.infra) continue;
                pages[site.url] = site;
            }
            resolve(pages);
        });
    });
}

function getHistory() {
    return new Promise(function(resolve, reject) {
        site.api.history({last: "2160h"}) // 90 days
            .get(function(err, data) {
            if (err) {
                reject(err);
                return;
            }

            let histories = {};

            for (let point of data) {
                if (!(point.url in histories)) {
                    histories[point.url] = [];
                }

                point.unix = Date.parse(point.timestamp);

                histories[point.url].push(point);
            }

            resolve(histories);
        });
    });
}