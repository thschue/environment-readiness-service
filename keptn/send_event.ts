// Description: This function triggers a github action
// Get Keptn context and variables

let params_raw = Deno.env.get("DATA");
let context_raw = Deno.env.get("CONTEXT");
let params;
let context;

if (params_raw != undefined) {
    params = JSON.parse(params_raw);
}

if (context_raw != undefined) {
    context = JSON.parse(context_raw);
}

let url = params.url;

let result = "my fancy result";

console.log(context);

let event = {
    "app": context.app,
    "workload": context.workload,
    "appVersion": context.appVersion,
    "workloadVersion": context.workloadVersion,
    "result": result,
    "eventType": context.eventType
}

let response = await fetch(url, {
    method: "POST",
    headers: {
        "Content-Type": "application/json",
    },
    body: JSON.stringify(event),
})

if (response.status == 200) {
    Deno.exit(0);
}
