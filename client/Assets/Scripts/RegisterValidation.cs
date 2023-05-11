using TMPro;
using UnityEngine;
using UnityEngine.Networking;

public class RegisterValidation : MonoBehaviour
{
    public TMP_InputField usernameRegister;
    public TMP_InputField emailRegister;
    public TMP_InputField passwordRegister;
    public TMP_InputField confirmPasswordRegister;

    public GameObject[] canvas1;
    public GameObject registerSuccessObject;
    public GameObject registerFailedObject;

    private void Start()
    {
        canvas1[0].SetActive(true);
    }

    public void CheckRegisterValidation()
    {
        string user = usernameRegister.text;
        string email = emailRegister.text;
        string pass = passwordRegister.text;
        string confirmPass = confirmPasswordRegister.text;

        if (string.IsNullOrEmpty(user) || string.IsNullOrEmpty(email) || string.IsNullOrEmpty(pass) || string.IsNullOrEmpty(confirmPass))
        {
            Debug.Log("Please enter your information.");
            return;
        }

        if (!IsValidEmail(email))
        {
            Debug.Log("Please enter a valid email.");
            return;
        }

        if (pass != confirmPass)
        {
            Debug.Log("Passwords do not match.");
            return;
        }

        Debug.Log("Button clicked with no parameters");
        // Create a request to the endpoint using POST method
        var req = new UnityWebRequest("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=AIzaSyBDEGHwraskC2U96zUEy5HwN3LFBlZNPDE", "POST");

        // Request body data
        var json = "{\"email\": \"" + email + "\", \"password\": \"" + pass + "\", \"returnSecureToken\": true}";
        byte[] jsonToSend = System.Text.Encoding.UTF8.GetBytes(json);

        // Set the request properties
        req.uploadHandler = new UploadHandlerRaw(jsonToSend);
        req.downloadHandler = new DownloadHandlerBuffer();
        req.SetRequestHeader("Content-Type", "application/json");

        // Send the request asynchronously 
        var request = req.SendWebRequest();

        while (!request.isDone)
        {
            // Waiting for the request to complete
        }

        if (req.result != UnityWebRequest.Result.Success)
        {
            Debug.Log("Form upload failed: " + req.error);
            // Handle error here
            registerFailedObject.SetActive(true);
        }
        else
        {
            Debug.Log("Form upload complete!");
            // Log the user information
            Debug.Log("Username: " + user);
            Debug.Log("Email: " + email);
            Debug.Log("Password: " + pass);
            Debug.Log("Confirm Password: " + confirmPass);

            // Handle the response from the Firebase Authentication API here
            registerSuccessObject.SetActive(true);
        }
    }

    private bool IsValidEmail(string email)
    {
        // You can implement your own email validation logic here
        // This is a basic example
        return email.Contains("@") && email.Contains(".");
    }
}
