## Keycloak Go OAuth2 `Client` Password Grant flow

This is an example golang application that interacts with a Keycloak instance as a `client` implementing OAuth2 password grant flow. The program is tested against `Keycloak 19.0.0`.

In order to run the application:

1. Set up a new realm (after login to your keycloak using admin) named `oauth-test` if it does not exist.

2. Setup your new application in keycloak : e.g. `http://localhost:8080/admin/master/console/#/<realm>/clients`. In `Client ID` field enter a `Client ID` (non empty ASCII text what will not be conflit with existing client settings). In the `Root URL` and `Home URL` field, enter `http://localhost:8082`, and in `Valid redirect URI` enter `/oauth/redirect`. Save and in `Credentials` tab you will see the generated credentials (including `client secret`).

3. Replace the values of the `ClientID`, `ClientSecret`, `token URL`, `auth URL` and `userinfo URL` in the `public/config.json` file which will be pickedup by both the front end (the htmls) and the backend (`main.go`). The values typically can be found from your Keycloak console page under `Realm settings/OpenID Endpoint Configuration` after you login to the console and select right Keycloak realm. Typically the [index.html] (which is under `public`) does not need to be changed. 

4. Start the server by executing `go run main.go`

5. Navigate to `http://localhost:8082` on your browser.

About the OAuth2 password grant flow:
1. If you correctly setup values in `public/config.json`, start the application. You will see a login page which ask you to enter user name and password. When you enter correctly (meaning a valid user exists in Keycloak in the given realm), the user name and password will be sent to bakcend (the `main.go` program), which will then send them to the Keycloak `token` endpoint together with `client ID and client secret, and grant_type=password` to get you authenticated and access token generated. 

2. The Keycloak `token` endpoint will respond with `access_token` and the backend program (`main.go`) will put you at `welcome` page. 

3. The `welcome` page will then ask the backend for `userinfo` and the backend will again ask Keycloak to provide userinfo by invoking the Keycloak `userinfo` endpoint. The userinfo (user name etc.) will then show on the welcome page.

* Note as the backend can record the username password and because there is no way to add other authentication mechanism into this flow, it is considered not secure and will be not be implemented in future versions of OAuth 2 compatible software.

