// Description: This function triggers a github action

// Get Keptn context and variables
let params_raw = Deno.env.get("DATA");
let context = Deno.env.get("CONTEXT");
let params;


if (params_raw != undefined) {
    params = JSON.parse(params_raw);
}

let url = params.url;

// Define the URL to the gitlab pipeline
// let url = params.url
let response = await fetch(url, {
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
            console.log("Infrastructure is not ready");
            Deno.exit(1);
        }
    });
