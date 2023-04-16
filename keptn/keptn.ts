// Description: This function triggers a github action
// Get Keptn context and variables
export function setContext() {
    let params_raw = Deno.env.get("DATA");
    let context_raw = Deno.env.get("CONTEXT");
    let params;
    let context;

    if (params_raw != undefined) {
        params = JSON.parse(params_raw);
    }

    if (context_raw != undefined) {
        context = JSON.parse(context_raw);
    } else {
        console.log("Context is not set - exiting")
        Deno.exit(1);
    }

    console.log(context);

    return [ params, context ];
}
export async function sendEvent(url: string, context: any, eventtype: string, result: string, success: boolean) {
    let event = {
        "app": context.app,
        "workload": context.workload,
        "appVersion": context.appVersion,
        "workloadVersion": context.workloadVersion,
        "result": result,
        "eventType": eventtype,
    }

    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(event),
    })

    console.log(JSON.stringify(event))
    if (success) {
        Deno.exit(0);
    }
    Deno.exit(1);
}