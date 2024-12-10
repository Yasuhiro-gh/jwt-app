# jwt-app
Simple API to generate and refresh pair <AccessToken, RefreshToken>
1. Build app
   ```
   docker compose build
   ```
2. Run app
   ```
   docker compose up
   ```

## API endpoints
### GET api/tokens/{uuid}
### POST api/refresh/refresh-data.json
---
### GET api/tokens/{uuid}
**Parameters**
| Name | Required| Type | Description                     |
|-----:|:-------:|:----:|:--------------------------------|
|`UserID`|required|UUID|User ID for which to refresh token|

**Response**

```
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM4MTI0MzksIlVzZXJJRCI6Ijg3NjRhNmZhLWIzM2EtNDZmNy1iYTFhLTZkOTE0MWQ3NTNmZiIsIklQQWRkciI6IjAuMC4wLjAiLCJTb21lUHJpdmFjeUluZm8iOiIyMDI0LTEyLTEwIDA2OjIzOjU5LjMxMzMyMDYyMSArMDAwMCBVVEMgbT0rMTAuMzg0MzYwOTA1In0.WJlBev5CreUrTBSt0UoB3LNs7AK_14p-PvZ5e8PpPafc8mQavXpkHbz_ceRIo0ekUu7EMyGvvGLY3UZwziPJhQ", // jwt
  "refresh_token": "MTczMzg0NzgzOQowLjAuMC4w" // base64
}
```
---

### POST api/refresh/refresh-data.json
**Parameters**
| Name | Required| Type | Description                        |
|-----:|:-------:|:----:|:-----------------------------------|
|`UserID`|required|UUID|User ID for which to refresh token   |
|`Refresh token`|required|String|Encoded base64 refresh token|

**Response**

```
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM4MTI0MzksIlVzZXJJRCI6Ijg3NjRhNmZhLWIzM2EtNDZmNy1iYTFhLTZkOTE0MWQ3NTNmZiIsIklQQWRkciI6IjAuMC4wLjAiLCJTb21lUHJpdmFjeUluZm8iOiIyMDI0LTEyLTEwIDA2OjIzOjU5LjMxMzMyMDYyMSArMDAwMCBVVEMgbT0rMTAuMzg0MzYwOTA1In0.WJlBev5CreUrTBSt0UoB3LNs7AK_14p-PvZ5e8PpPafc8mQavXpkHbz_ceRIo0ekUu7EMyGvvGLY3UZwziPJhQ", // jwt
  "refresh_token": "MTczMzg0NzgzOQowLjAuMC4w" // base64
}

