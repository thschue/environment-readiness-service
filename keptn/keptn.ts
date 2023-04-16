// Description: This function triggers a github action
// Get Keptn context and variables
export function setContext(): [any, any] {
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
    return params, context;
}
export function sendEvent(url: string, context: any, result: string) {
    let event = {
        "app": context.app,
        "workload": context.workload,
        "appVersion": context.appVersion,
        "workloadVersion": context.workloadVersion,
        "result": result,
        "eventType": context.eventType
    }

    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(event),
    })
}
