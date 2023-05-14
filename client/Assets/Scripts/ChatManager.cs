using System;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

public class ChatManager : MonoBehaviour
{
    public TMP_InputField messageInputField;
    public Button sendButton;
    public GameObject messageContainerPrefab;
    private GameObject messageContainer;
    private SocketClient socketClient;

    private void Start()
    {
        sendButton.onClick.AddListener(OnSendButtonClick);
        messageInputField.onEndEdit.AddListener(OnEndEdit);
        CreateMessageContainer();

        socketClient = SocketClient.Instance; // Get reference to SocketClient instance

        // Initialize SocketClient if it's not already initialized
        if (!socketClient.IsInitialized)
        {
            // Initialize SocketClient with wss:// host and random userId
            string host = "wss://werewolf-app.fly.dev/";
            string userId = Guid.NewGuid().ToString(); // Generate a random userId using GUID
            socketClient.Initialize(host, userId);
            socketClient.Connect();
        }
        else
        {
            Debug.LogWarning("SocketClient is already initialized!");
        }
    }

    private void OnDestroy()
    {
        socketClient.Disconnect();
    }

    private void OnSendButtonClick()
    {
        SendMessage();
    }

    private void OnEndEdit(string text)
    {
        if (Input.GetKey(KeyCode.Return) || Input.GetKey(KeyCode.KeypadEnter))
        {
            SendMessage();
        }
    }

    private void SendMessage()
    {
        string message = messageInputField.text;
        socketClient.ChatInRoom(message);

        // Clear the input field after sending the message
        messageInputField.text = "";
    }

    public void ReceiveMessageFromServer(string message)
    {
        if (messageContainer != null)
        {
            // Display the received message in the message container
            messageContainer.GetComponent<MessageContainerScript>().AddMessage("Server: " + message);
        }
        else
        {
            Debug.LogWarning("Message container is missing!");
        }
    }

    private void CreateMessageContainer()
    {
        if (messageContainerPrefab != null)
        {
            messageContainer = Instantiate(messageContainerPrefab, transform);
        }
        else
        {
            Debug.LogWarning("Message container prefab is missing!");
        }
    }
}
