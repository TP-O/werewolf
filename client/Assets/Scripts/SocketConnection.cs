using System;
using System.Collections;
using System.Collections.Generic;
using SocketIOClient;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class SocketConnection : MonoBehaviour
{
    public SocketIOUnity socket;
    public TMP_Text receivedText;

    private string token;

    void Start()
    {
        // Check if the URI is valid
        var uri = new Uri("https://werewolf-app.fly.dev/");
        socket = new SocketIOUnity(uri, new SocketIOOptions
        {
            Query = new Dictionary<string, string>
            {
                {"token", "UNITY" }
            },
            Transport = SocketIOClient.Transport.TransportProtocol.WebSocket
        });

        // Event listeners for Socket.IO events
        socket.OnConnected += (sender, e) =>
        {
            Debug.Log("Socket connected.");
            receivedText.text += "Socket connected.\n";
        };

        socket.OnDisconnected += (sender, e) =>
        {
            if (e == "io server disconnect")
            {
                // Thực hiện xử lý khi máy chủ Socket.IO ngắt kết nối
                // Ví dụ: Thử kết nối lại sau một khoảng thời gian
                StartCoroutine(ReconnectAfterDelay(5f));
            }
            else
            {
                // Xử lý các trường hợp ngắt kết nối khác
            }

            receivedText.text += "Socket disconnected. Reason: " + e + "\n";
        };

        // Event listeners for custom events
        socket.On("eventName", response =>
        {
            var eventData = response.GetValue<string>();
            // Do something with eventData
            Debug.Log("Received Event: " + eventData);
            receivedText.text += "Event Data: " + eventData + "\n";
        });

        Debug.Log("Connecting to socket...");
        receivedText.text += "Connecting to socket...\n";
        socket.Connect();

        socket.On("Error", response =>
        {
            var errorData = response.GetValue<string>();
            // Do something with errorData
            Debug.Log("Received Error: " + errorData);
            receivedText.text += "Error: " + errorData + "\n";
        });

        socket.On("Success", response =>
        {
            var successData = response.GetValue<string>();
            // Do something with successData
            Debug.Log("Received Success: " + successData);
            receivedText.text += "Success: " + successData + "\n";
        });

        socket.On("RoomMessage", response =>
        {
            var roomMessageData = response.GetValue<string>();
            // Do something with roomMessageData
            Debug.Log("Received RoomMessage: " + roomMessageData);
            receivedText.text += "RoomMessage: " + roomMessageData + "\n";
        });

        socket.On("RoomChange", response =>
        {
            var roomChangeData = response.GetValue<string>();
            // Do something with roomChangeData
            Debug.Log("Received RoomChange: " + roomChangeData);
            receivedText.text += "RoomChange: " + roomChangeData + "\n";
        });

    }

    public void EmitTest(string eventName, string data)
    {
        socket.Emit(eventName, data);
    }

    private string GetToken()
    {
        if (string.IsNullOrEmpty(token))
        {
            token = LoginValidation.GetToken();
        }
        return token;
    }

    private IEnumerator ReconnectAfterDelay(float delay)
    {
        yield return new WaitForSeconds(delay);
        socket.Connect();
    }

}
