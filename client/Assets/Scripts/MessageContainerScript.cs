using TMPro;
using UnityEngine;

public class MessageContainerScript : MonoBehaviour
{
    public TMP_Text messageText;

    private void Start()
    {
        ClearMessages();
    }

    public void AddMessage(string message)
    {
        // Append the message to the message text
        messageText.text += message + "\n";
    }

    public void ClearMessages()
    {
        // Clear the message text
        messageText.text = "";
    }
}
