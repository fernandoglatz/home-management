<!-- Generator: Widdershins v4.0.1 -->

<h1 id="home-management">home-management v1.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api">/api</a>

License: <a href="http://www.apache.org/licenses/LICENSE-2.0.html">Apache 2.0</a>

# Authentication

- HTTP Authentication, scheme: basic 

* API Key (Bearer)
    - Parameter Name: **X-AUTHORIZATION**, in: header. Generated by /authentication

<h1 id="home-management-event">event</h1>

## get__event

`GET /event`

*Get events*

> Example responses

> 200 Response

```json
[
  {
    "createdAt": "string",
    "device": "string",
    "home": "string",
    "id": "string",
    "type": "string",
    "updatedAt": "string",
    "version": "string"
  }
]
```

<h3 id="get__event-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<h3 id="get__event-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[entity.Event](#schemaentity.event)]|false|none|none|
|» createdAt|string|false|none|none|
|» device|string|false|none|none|
|» home|string|false|none|none|
|» id|string|false|none|none|
|» type|string|false|none|none|
|» updatedAt|string|false|none|none|
|» version|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## put__event

`PUT /event`

*Create event*

> Body parameter

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="put__event-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[entity.Event](#schemaentity.event)|true|body|

> Example responses

> 200 Response

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="put__event-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[entity.Event](#schemaentity.event)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

## get__event_{id}

`GET /event/{id}`

*Get event*

<h3 id="get__event_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|id|

> Example responses

> 200 Response

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="get__event_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[entity.Event](#schemaentity.event)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

## put__event_{id}

`PUT /event/{id}`

*Update event*

> Body parameter

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="put__event_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|id|
|body|body|[entity.Event](#schemaentity.event)|true|body|

> Example responses

> 200 Response

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="put__event_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[entity.Event](#schemaentity.event)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

## post__event_{id}

`POST /event/{id}`

*Update event*

> Body parameter

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="post__event_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|id|
|body|body|[entity.Event](#schemaentity.event)|true|body|

> Example responses

> 200 Response

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="post__event_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[entity.Event](#schemaentity.event)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

## delete__event_{id}

`DELETE /event/{id}`

*Delete event*

<h3 id="delete__event_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|id|

> Example responses

> 400 Response

```json
{
  "code": "string",
  "message": "string"
}
```

<h3 id="delete__event_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

## patch__event_{id}

`PATCH /event/{id}`

*Update event*

> Body parameter

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="patch__event_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|id|
|body|body|[entity.Event](#schemaentity.event)|true|body|

> Example responses

> 200 Response

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}
```

<h3 id="patch__event_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[entity.Event](#schemaentity.event)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[response.Response](#schemaresponse.response)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[response.Response](#schemaresponse.response)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="home-management-health">health</h1>

## get__health

`GET /health`

*Get health*

> Example responses

> 200 Response

```
"string"
```

<h3 id="get__health-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|string|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|string|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|string|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="home-management-events">events</h1>

## head__v1_events_{id}

`HEAD /v1/events/{id}`

*Check if event exists*

<h3 id="head__v1_events_{id}-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|Event ID|

<h3 id="head__v1_events_{id}-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not Found|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_entity.Event">entity.Event</h2>
<!-- backwards compatibility -->
<a id="schemaentity.event"></a>
<a id="schema_entity.Event"></a>
<a id="tocSentity.event"></a>
<a id="tocsentity.event"></a>

```json
{
  "createdAt": "string",
  "device": "string",
  "home": "string",
  "id": "string",
  "type": "string",
  "updatedAt": "string",
  "version": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|createdAt|string|false|none|none|
|device|string|false|none|none|
|home|string|false|none|none|
|id|string|false|none|none|
|type|string|false|none|none|
|updatedAt|string|false|none|none|
|version|string|false|none|none|

<h2 id="tocS_response.Response">response.Response</h2>
<!-- backwards compatibility -->
<a id="schemaresponse.response"></a>
<a id="schema_response.Response"></a>
<a id="tocSresponse.response"></a>
<a id="tocsresponse.response"></a>

```json
{
  "code": "string",
  "message": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string|false|none|none|
|message|string|false|none|none|

