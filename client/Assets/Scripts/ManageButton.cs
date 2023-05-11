using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;

public class ManageButton : MonoBehaviour
{
    [SerializeField] private Button btn = null;

    private void Awake()
    {
        // adding a delegate with no parameters
        btn.onClick.AddListener(NoParameterOnClick);

        // adding a delegate with parameters
        btn.onClick.AddListener(delegate { ParameterOnClick("Button was pressed!"); });
    }

    private void NoParameterOnClick()
    {
        AAaa();
    }

    public void AAaa()
    {
        Debug.Log("Button clicked with no parameters");

        var req = UnityWebRequest.Post("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=AIzaSyBDEGHwraskC2U96zUEy5HwN3LFBlZNPDE", "");
        req.SetRequestHeader("Content-Type", "application/json");

        string json = "{\"email\": \"69phongle@gmail.com\",\"password\": \"dsdsdsdsdsdsdsds\", \"returnSecureToken\": true}";
        byte[] jsonBytes = System.Text.Encoding.UTF8.GetBytes(json);
        req.uploadHandler = new UploadHandlerRaw(jsonBytes);
        req.downloadHandler = new DownloadHandlerBuffer();

        var operation = req.SendWebRequest();

        operation.completed += asyncOperation =>
        {
            var request = operation.webRequest;

            if (request.result != UnityWebRequest.Result.Success)
            {
                Debug.Log(request.error);
            }
            else
            {
                Debug.Log("Form upload complete!");
                Debug.Log("Response Code: " + request.responseCode);
                Debug.Log("Response Result: " + request.downloadHandler.text);
            }
        };
    }

    private void ParameterOnClick(string test)
    {
        Debug.Log(test);
    }
}
