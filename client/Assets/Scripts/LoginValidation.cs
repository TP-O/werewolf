using UnityEngine;
using UnityEngine.Networking;
using TMPro;
using System.Text;
using System.Collections;

public class LoginValidation : MonoBehaviour
{
    public TMP_InputField username;
    public TMP_InputField password;

    public GameObject[] canvas;
    public GameObject loginFailedObject;

    public LoginMenu loginMenu;
    private static string token;

    private void Start()
    {
        canvas[0].SetActive(true);
    }

    public void CheckLoginValidation()
    {
        string email = username.text;
        string pass = password.text;

        Debug.Log("Button clicked with no parameters");
        // Create a request to the endpoint using the POST method
        var req = UnityWebRequest.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyBDEGHwraskC2U96zUEy5HwN3LFBlZNPDE", "");

        // Request data
        var json = "{\"email\": \"" + email + "\", \"password\": \"" + pass + "\", \"returnSecureToken\": true}";
        byte[] jsonToSend = Encoding.UTF8.GetBytes(json);

        // Set request properties
        req.uploadHandler = new UploadHandlerRaw(jsonToSend);
        req.downloadHandler = new DownloadHandlerBuffer();
        req.SetRequestHeader("Content-Type", "application/json");

        // Send the request asynchronously
        var request = req.SendWebRequest();

        // Coroutine to handle the response
        StartCoroutine(HandleLoginResponse(request));
    }

    private IEnumerator HandleLoginResponse(UnityWebRequestAsyncOperation asyncOperation)
    {
        var request = asyncOperation.webRequest;

        yield return asyncOperation;

        if (request.result != UnityWebRequest.Result.Success)
        {
            Debug.Log("Form upload failed: " + request.error);
            // Handle error here
            loginFailedObject.SetActive(true);
        }
        else
        {
            Debug.Log("Form upload complete!");

            // Get the response data as a JSON string
            string responseText = request.downloadHandler.text;
            Debug.Log("Response: " + responseText);

            // Parse the JSON response
            // You can use a JSON parsing library like JsonUtility or Newtonsoft.Json
            // to deserialize the JSON string into an object and extract the token

            // Example using JsonUtility (assuming the response has a "token" field)
            LoginResponseData responseData = JsonUtility.FromJson<LoginResponseData>(responseText);
            token = responseData.token;

            // Print the token
            Debug.Log("Token: " + token);

            // Call the NextScene() function in LoginMenu script
            loginMenu.NextScene();
        }

        // Dispose the UnityWebRequest to prevent memory leak
        request.Dispose();
    }

    public static string GetToken()
    {
        return token; // Phương thức để truy cập token từ bên ngoài lớp
    }

}
