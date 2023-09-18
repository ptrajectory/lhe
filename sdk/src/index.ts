import NetClient from "../lib/net";
import { EmailEvent, WebhookEvent, tEmailEvent, tWebhookEvent } from "../lib/zod";


class LHE {

    private client

    constructor(serverUrl: string, apiKey: string){
        const baseHeaders = new Headers()
        baseHeaders.append("Authorization", `Bearer ${apiKey}`)
        baseHeaders.append("Content-Type", "application/json")
        this.client = new NetClient(serverUrl, baseHeaders)
    }


    async createEmail(data: tEmailEvent){

        const parsed = EmailEvent.safeParse(data)
        
        if(!parsed.success) throw new Error(JSON.stringify(parsed.error.formErrors.fieldErrors))

        const parsedData = parsed.data

        const resp = await this.client.post("/events/emails", {
            body: parsedData
        })


        if(resp.ok) return resp.data 

        throw new Error(resp.statusText)

    }


    async createWebhook(data: tWebhookEvent){

        const parsed = WebhookEvent.safeParse(data)

        if(!parsed.success) throw new Error(JSON.stringify(parsed.error.formErrors.fieldErrors))

        const parsedData = parsed.data

        const resp = await this.client.post("/events/webhooks", {
            body: parsedData
        })

        if(resp.ok) return resp.data 

        throw new Error(resp.statusText)

    }

}


export default function createLHEClient(url: string, apiKey: string){

    const lhe = new LHE(url,apiKey)

    return lhe

}