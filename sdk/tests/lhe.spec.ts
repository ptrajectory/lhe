import createLHEClient from "../src"


const client = createLHEClient(process.env.API_URL as string, process.env.API_KEY as string)


describe("TEST EMAIL AND WEBHOOK CREATION",()=>{

    it("Should add an email", (done)=>{

        client.createEmail({
            payload: {
                body: "Test from api",
                from: process.env.SENDER_EMAIL as string,
                subject: "THIS IS JUST A TEST"
            },
            receiver: {
                email: process.env.RECEIVER_EMAIL as string
            },
            sender: {
                email: process.env.SENDER_EMAIL as string,
                id: "TEST1234",
                name: "Duck"
            }
        })
        .then((data)=>{
            console.log("SOME RESPONSE::",data)
            done()
        })
        .catch((e)=>{
            console.log("Something went wrong::",e)
            done(e)
        })

    })


    it("Should add a webhook", (done)=>{

        client.createWebhook({
            payload: {
                resource: "order",
                type: "order.created",
                object: {
                    order_id: "123"
                }
            },
            receiver: {
                url: process.env.TEST_WEBHOOK_ENDPOINT as string, 
                headers: {}
            },
            sender: {
                name: "My APP",
                id: "app_id",
                email: "email@email.com"
            },
            type: "order_created"
        }).then((data)=>{
            console.log("SOME RESPONSE::",data)
            done()
        })
        .catch((e)=>{
            console.log("Something went wrong::",e)
            done(e)
        })

    })
    

})