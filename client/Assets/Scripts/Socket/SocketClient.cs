using System;
using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;

public class SocketClient : MonoBehaviour
{
    public bool IsConnecting { get; private set; } = false;
    public bool IsConnected { get; private set; } = false;
    public bool HasRoom => !string.IsNullOrEmpty(RoomId);
    public string RoomId { get; private set; } = "";
    public string[] Chats => _chats.ToArray();
    public string UserId { get; private set; }
    public string ApiKey { get; set; }

    public SocketClient socketClient;

    private WebSocket _socket;
    private string _uri;
    private string _userId;
    private List<string> _chats = new List<string>();
    public event Action<string> OnEventReceived;

    public void Initialize(string host, string userId)
    {
        _uri = $"{host}?userId={userId}";
        _userId = userId;
        Debug.LogFormat("SocketClient. userId: {0}, uri: {1}", userId, _uri);
    }

    public void Initialize(string host, string userId, string apiKey)
    {
        // Các dòng code khác trong phương thức Initialize...

        ApiKey = apiKey;
    }


    public void Connect()
    {
        if (!IsConnected && !IsConnecting)
        {
            Debug.Log("Connect");
            _socket = new WebSocket(_uri);

            // Standard events
            _socket.OnOpen += OnConnect;
            _socket.OnError += OnError;
            _socket.OnClose += OnDisconnect;

            // Room events
            _socket.OnMessage += OnMessageReceived;

            _socket.ConnectAsync();
            IsConnecting = true;
        }
    }

    public void Disconnect()
    {
        if (IsConnecting || IsConnected)
        {
            Debug.Log("Disconnect");

            if (_socket != null)
            {
                _socket.OnOpen -= OnConnect;
                _socket.OnError -= OnError;
                _socket.OnClose -= OnDisconnect;
                _socket.OnMessage -= OnMessageReceived;

                _socket.Close();
                _socket = null;
            }

            RoomId = "";
            IsConnecting = false;
            IsConnected = false;
        }
    }

    public void CreateRoom()
    {
        if (_socket != null && !HasRoom)
        {
            Debug.Log("CreateRoom");
            _socket.Send("createRoom");
        }
    }

    public void GetRoom()
    {
        if (_socket != null)
        {
            Debug.Log("GetRoom");
            _socket.Send("getRoom");
        }
    }

    public void LeaveRoom()
    {
        if (HasRoom && _socket != null)
        {
            Debug.Log("LeaveRoom");
            _socket.Send("leaveRoom");
        }
    }

    public void FindRoom()
    {
        if (!HasRoom && _socket != null)
        {
            Debug.Log("FindRoom");
            _socket.Send("findRoom ");
        }
    }

    public void JoinRoom(string roomId)
    {
        if (!HasRoom && _socket != null)
        {
            Debug.Log("JoinRoom");
            _socket.Send("joinRoom " + roomId);
        }
    }

    public void ChatInRoom(string message)
    {
        if (HasRoom && _socket != null)
        {
            Debug.Log("ChatInRoom: " + message);
            _socket.Send("chatInRoom " + message);
        }
    }

    private void OnDestroy()
    {
        Disconnect();
    }

    private void OnConnect(object sender, EventArgs e)
    {
        Debug.Log("Socket connected.");
        IsConnecting = false;
        IsConnected = true;
    }

    private void OnMessageReceived(object sender, MessageEventArgs e)
    {
        string message = e.Data;
        Debug.Log("Received message: " + message);

        // Xử lý dữ liệu nhận được từ server trước khi chạy code
        // Ví dụ: in ra dữ liệu sự kiện trả về
        OnEventReceived?.Invoke(message);
    }

    private void OnDisconnect(object sender, CloseEventArgs e)
    {
        Debug.Log("_onDisconnect");
        Disconnect();
    }


    private void OnError(object sender, ErrorEventArgs e)
    {
        Debug.LogError("Error: " + e.Message);
        Disconnect();
    }

    private static SocketClient _instance;
    internal bool IsInitialized;

    public static SocketClient Instance
    {
        get
        {
            if (_instance == null)
            {
                _instance = FindObjectOfType<SocketClient>();

                if (_instance == null)
                {
                    GameObject obj = new GameObject("SocketClient");
                    _instance = obj.AddComponent<SocketClient>();
                }
            }

            return _instance;
        }
    }
}

