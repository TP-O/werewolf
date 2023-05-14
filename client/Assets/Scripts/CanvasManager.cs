using UnityEngine;
using UnityEngine.UI;

public class CanvasManager : MonoBehaviour
{
    public GameObject[] buttonObjects; // Array containing the Button GameObjects
    private SocketClient socketClient;


    private void Start()
    {
        foreach (var buttonObject in buttonObjects)
        {
            Button button = buttonObject.GetComponent<Button>();

            if (button != null)
            {
                if (buttonObject.name == "Create Game")
                {
                    button.onClick.AddListener(OnCreateRoomButtonClicked);
                }
                else if (buttonObject.name == "Find Game")
                {
                    button.onClick.AddListener(OnFindGameButtonClicked);
                }
                // Add other button event handling cases if needed
            }
        }

        //socketClient = socketClient;
        socketClient = FindObjectOfType<SocketClient>();

        if (socketClient == null)
        {
            Debug.LogError("SocketClient component not found in the scene.");
            return;
        }

        string userId = System.Guid.NewGuid().ToString();
        string apiKey = "AIzaSyBDEGHwraskC2U96zUEy5HwN3LFBlZNPDE"; // Thay YOUR_API_KEY bằng giá trị thực của API key
        socketClient.Initialize("wss://werewolf-app.fly.dev/", userId, apiKey);
        socketClient.Connect();
    }

    private void OnDestroy()
    {
        if (socketClient != null)
        {
            socketClient.OnEventReceived -= OnServerEventReceived;
        }
    }

    public void OnCreateRoomButtonClicked()
    {
        if (socketClient != null)
        {
            socketClient.CreateRoom();
        }
        else
        {
            Debug.LogError("SocketClient is not assigned.");
        }
    }

    public void OnFindGameButtonClicked()
    {
        if (socketClient != null)
        {
            socketClient.FindRoom();
        }
        else
        {
            Debug.LogError("SocketClient is not assigned.");
        }
    }

    private void OnServerEventReceived(string eventData)
    {
        Debug.Log("Server event received: " + eventData);
        // Handle server data if needed
    }
}
