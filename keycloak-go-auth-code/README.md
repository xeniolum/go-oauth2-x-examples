## Keycloak Go OAuth2 Auth-code Flow Example

This is an example golang application that interacts with a Keycloak instance by implementing OAuth2 auth-code grant flow. The program is tested against `Keycloak 19.0.0`.

In order to run the application:

1. Set up a new realm (after login to your keycloak using admin) named `oauth-test` if it does not exist.

2. Setup your new application in keycloak : e.g. `http://localhost:8080/admin/master/console/#/<realm>/clients`. In `Client ID` field enter a <Client ID> (non empty ASCII text what will not be conflit with existing client settings). In the `Root URL` and `Home URL` field, enter `http://localhost:8082`, and in `Valid redirect URI` enter `/oauth/redirect`. Save and in `Credentials` tab you will see the generated credentials (including `client secret`).

3. Replace the values of the <ClientID>, <ClientSecret>, <token URL>, <auth URL> and <userinfo URL> in the `public/config.json` file which will be pickedup by both the front end (the htmls) and the backend (`main.go`). The values typically can be found from your Keycloak console page under `Realm settings/OpenID Endpoint Configuration` after you login to the console and select right Keycloak realm. Typically the [index.html] (which is under `public`) does not need to be changed. 

4. Start the server by executing `go run main.go`

5. Navigate to `http://localhost:8082` on your browser.

About the OAuth2 auth-code grant flow (which is the most important OAuth2 flow):
1. Start the application, you will see a link which says: `Login with keycloak` and when you click on that link and if you correctly setup values in `public/config.json`, the link will bring you to Keycloak realm login page. (This step your Client ID and Client secret will be sent to Keycloak.)

2. If you successfully login, your browser will be redirected to `/oauth/redirect` with generated `auth code` sent back.

3. The backend `main.go` will process the code and require token by invoking `token URL` (this process is called `CODE FOR TOKEN`), then it will forward you to the `welcome.html` (whihc is under `public` folder) page with `access_token` attached in the URL, and by now you should be sure that you are authorized to use the client and you are authenticated.

4. The `welcome.html` page will typically make another request to ask the server for user information (which in our case the user name) using `/oauth/usreinfo` and for that purpose the server in background will be invoking Keycloak `userinfo URL` endpoint as setup in `public/config.json` file.
