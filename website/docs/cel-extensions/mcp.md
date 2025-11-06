# MCP Library

The MCP (Match Condition Parameters) CEL extension enables parsing and extracting arguments from an MCP JSON payload, as specified in [Envoy's Match Condition Parameters](https://github.com/envoyproxy/envoy/blob/main/api-docs/xds_protocol.rst#match-condition-parameters).

This library allows Kyverno policies to extract, filter, and use MCP/proxy metadata in CEL expressions.

---

## Types

### `<MCP>`

*CEL Type / Proto*: `mcp.MCP`

This is an opaque type that represents an MCP parser instance.

### `<MCPRequest>`

*CEL Type / Proto*: `mcp.MCPRequest`

Represents a parsed MCP request, providing accessors for extracting arguments.

The `MCPRequest` struct contains the parsed MCP request data. It includes:

- **`Method`** (string): The MCP method name
- **`ID`** (RequestId): The request identifier
- **`Paginated`** (*PaginatedParams): Pagination parameters for all list methods

**Tool-related fields:**
- **`ToolCall`** (*CallToolParams): Tool call parameters

**Resource-related fields:**
- **`ResourceRead`** (*ReadResourceParams): Resource read parameters
- **`ResourceSubscribe`** (*SubscribeParams): Resource subscription parameters
- **`ResourceUnsubscribe`** (*UnsubscribeParams): Resource unsubscription parameters

**Prompt-related fields:**
- **`PromptGet`** (*GetPromptParams): Prompt retrieval parameters

**Lifecycle and utility fields:**
- **`Initialize`** (*InitializeParams): Lifecycle initialization parameters
- **`CreateMessage`** (*CreateMessageParams): Message creation parameters
- **`Elicitation`** (*ElicitationParams): Elicitation parameters
- **`Complete`** (*CompleteParams): Completion utility parameters
- **`SetLogLevel`** (*SetLevelParams): Log level setting parameters

These fields are conditionally populated based on the request method and type that was received (eg: ToolCall will only be populated if MCP method was `tools/call`).

---

## Constants

The following values representing known MCP methods are available as constants for use in expressions (exposed as `mcp.<name>`):

- `mcp.InitializeMethod`
- `mcp.PingMethod`
- `mcp.ResourcesListMethod`
- `mcp.ResourcesTemplatesListMethod`
- `mcp.ResourcesReadMethod`
- `mcp.PromptsListMethod`
- `mcp.PromptsGetMethod`
- `mcp.ToolsListMethod`
- `mcp.ToolsCallMethod`
- `mcp.SetLogLevelMethod`
- `mcp.ElicitationCreateMethod`

---

## Functions

### Parse

Parses a raw MCP JSON string using a given MCP parser instance.

#### Signature

```
MCP.Parse(<string> value) -> MCPRequest
```

#### Example

```cel
// "rawMcpJson" holds the MCP JSON string (usually the proxy's metadata)
mcpRequest := mcpInstance.Parse(rawMcpJson)
```

---

### GetStringArgument

Returns the string value of an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetStringArgument(<string> key, <string> default) -> <string>
```

#### Example

```cel
userRole := mcpRequest.GetStringArgument("role", "guest")
```

---

### GetIntArgument

Returns the integer value of an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetIntArgument(<string> key, <int> default) -> <int>
```

#### Example

```cel
requestCount := mcpRequest.GetIntArgument("count", 0)
```

---

### GetFloatArgument

Returns the float (double) value of an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetFloatArgument(<string> key, <double> default) -> <double>
```

#### Example

```cel
confidence := mcpRequest.GetFloatArgument("confidence", 0.0)
```

---

### GetBoolArgument

Returns the boolean value of an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetBoolArgument(<string> key, <bool> default) -> <bool>
```

#### Example

```cel
isAdmin := mcpRequest.GetBoolArgument("admin", false)
```

---

### GetStringSliceArgument

Returns a list of strings for an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetStringSliceArgument(<string> key, <list<string>> default) -> <list<string>>
```

#### Example

```cel
groups := mcpRequest.GetStringSliceArgument("groups", ["public"])
```

---

### GetIntSliceArgument

Returns a list of integers for an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetIntSliceArgument(<string> key, <list<int>> default) -> <list<int>>
```

#### Example

```cel
ids := mcpRequest.GetIntSliceArgument("ids", [])
```

---

### GetFloatSliceArgument

Returns a list of floats for an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetFloatSliceArgument(<string> key, <list<double>> default) -> <list<double>>
```

---

### GetBoolSliceArgument

Returns a list of booleans for an argument given its key, or a default if missing.

#### Signature

```
MCPRequest.GetBoolSliceArgument(<string> key, <list<bool>> default) -> <list<bool>>
```

---

## Example: Using MCP with Kyverno CEL

Let's say we want to check a tool's usage:

```cel
// Suppose you have received an MCP JSON string in input.mcp
mcpReq = mcpInstance.Parse(input.mcp)
// Check the MCP request's method & tool usage name
isAllowed = mcpReq.Method == mcp.ToolsCallMethod && mcpReq.ToolCall.Name == "shell"
// Check for a particular tool argument
command = mcpReq.GetStringArgument("command", "")
isAuthorized = isCallMethod && isShellTool && command in ["kubectl", "docker", ""]
```

---

