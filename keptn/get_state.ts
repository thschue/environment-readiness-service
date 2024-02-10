// Description: This function triggers a github action
// Get Keptn context and variables
import { sleep } from "https://deno.land/x/sleep/mod.ts";

let params_raw = Deno.env.get("DATA");
let context = Deno.env.get("KEPTN_CONTEXT");
let params;

if (params_raw != undefined) {
    params = JSON.parse(params_raw);
}

let url = params.url;
let timeout = params.timeout;
let retries = params.retries;

if (timeout == undefined) {
    timeout = 5;
}

if (retries == undefined) {
    retries = 10;
}

let response;

let tries = 0;
while (tries < retries) (
    response = await fetch(url, {
        method: 'GET',
        headers: {
            Accept: 'application/json',
        },
    })
        .then(function (response) {
            // The response is a Response instance.
            // You parse the data into a useable format using `.json()`
            return response.json();
        }).then(function (data) {
            // `data` is the parsed version of the JSON returned from the above endpoint.
            if (data.infraReady) {
                console.log("Infrastructure is ready");
                Deno.exit(0);
            } else {
                console.log("Infrastructure is not ready - retrying");
                tries++;
                if (tries >= retries) {
                    console.log("Infrastructure is not ready after " + retries + " retries - exiting");
                    Deno.exit(1);
                }
            }
        }),
    await sleep(timeout)
)
