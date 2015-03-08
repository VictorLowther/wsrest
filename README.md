wsrest is a sorta RESTy wsman gateway.

Once built (with go build -i), just run it with a --listen argument:

    ./wsrest -listen=127.0.0.1:31337

The server will start listening for incoming POSTs to /.  The posts must
contain a JSON payload:

    {
        "Endpoint": "https://the.WSMAN.endpoint:443/wsman",
        "Username": "username",
        "Password": "password",
        "Method":   "http://the.wsman.method",
        "ResourceURI: "http://the.resource/to/use",
        "Options": ["key1", "value1", "key2", "value2"],
        "Selectors": ["name1", "selector1", "name2", "selector2"],
        "Parameters": ["name1", "parameter1", "name2", "parameter2"]
    }

Method, ResourceURI, Options, Selectors, and Parameters will be used to
build a WSMAN-compliant SOAP message, which will be sent to Endpoint using
Username and Password for authentication.  The XML response (if any) will be
returned unaltered.

In common cases, method will be one of:

* "Identify", which will perform WSMAN Identification against the endpoint.
* "Get", which will perform a WSMAN Get of ResourceURI
* "Enumerate", which will perform WSMAN Enumerate of ResourceURI
* "EnumerateEPR", which will enumerate the endpoint resources of ResourceURI
* "Invoke", which will split ResourceURI into a resource / method pair at the last slash
  and invohe method on resource.

Anything else will be invoked directly.


For example:

    $ curl -X POST http://127.0.0.1:18376 -d '{
        "Endpoint": "https://192.168.128.22:443/wsman",
        "Username": "root",
        "Password": "password",
        "Method": "Invoke",
        "ResourceURI": "http://schemas.dell.com/wbem/wscim/1/cim-schema/2/DCIM_CSPowerManagementService/RequestPowerStateChange",
        "Selectors": ["Name", "pwrmgtsvc:1",
            "CreationClassName", "DCIM_CSPowerManagementService",
            "SystemCreationClassName", "DCIM_SPComputerSystem",
            "SystemName", "systemmc"],
            "Parameters": ["PowerState","2"]}'

will turn on a Dell system, and give back the following response:

    <?xml version="1.0" encoding="UTF-8"?>
    <s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
        xmlns:wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing"
        xmlns:n1="http://schemas.dell.com/wbem/wscim/1/cim-schema/2/DCIM_CSPowerManagementService">
       <s:Header>
         <wsa:To>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</wsa:To>
         <wsa:Action>http://schemas.dell.com/wbem/wscim/1/cim-schema/2/DCIM_CSPowerManagementService/RequestPowerStateChangeResponse</wsa:Action>
         <wsa:RelatesTo>uuid:3c1c6d34-a13a-46e5-a325-ba3362f46ce0</wsa:RelatesTo>
         <wsa:MessageID>uuid:065558d8-10d3-10d3-8041-76c44112bcf8</wsa:MessageID>
       </s:Header>
       <s:Body>
         <n1:RequestPowerStateChange_OUTPUT>
           <n1:ReturnValue>0</n1:ReturnValue>
         </n1:RequestPowerStateChange_OUTPUT>
      </s:Body>
    </s:Envelope>
