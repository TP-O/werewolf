using UnityEngine;
using UnityEngine.UI;

public class ButtonHandler : MonoBehaviour
{
    private SocketClient socketClient;

    private void Start()
    {
        // Kiểm tra và lấy tham chiếu đến SocketClient
        socketClient = FindObjectOfType<SocketClient>();

        if (socketClient == null)
        {
            Debug.LogError("SocketClient component not found in the scene.");
        }

        // Gắn phương thức xử lý sự kiện cho nút "Create Room"
        Button createRoomButton = GameObject.Find("Create Game").GetComponent<Button>();
        createRoomButton.onClick.AddListener(OnCreateRoomButtonClicked);

        // Gắn phương thức xử lý sự kiện cho nút "Find Game"
        Button findGameButton = GameObject.Find("Find Game").GetComponent<Button>();
        findGameButton.onClick.AddListener(OnFindGameButtonClicked);
    }

    private void OnCreateRoomButtonClicked()
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

    private void OnFindGameButtonClicked()
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
}
