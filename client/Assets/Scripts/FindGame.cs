using System.Collections;
using UnityEngine;
using UnityEngine.Networking;
using UnityEngine.SceneManagement;

public class FindGame : MonoBehaviour
{
    private string serverURL = "https://werewolf-app.fly.dev"; // URL của máy chủ
    public GameObject createFailedObject;
    public GameObject[] objectsToDisable;

    public void StartFindingGame()
    {
        StartCoroutine(SendRequestToServer());
    }

    private IEnumerator SendRequestToServer()
    {
        // Gửi yêu cầu tới máy chủ
        using (UnityWebRequest request = UnityWebRequest.Get(serverURL))
        {
            yield return request.SendWebRequest();

            if (request.result == UnityWebRequest.Result.Success)
            {
                // Xử lý phản hồi từ máy chủ
                string responseText = request.downloadHandler.text;
                Debug.Log("Response: " + responseText);

                // Enable UI thông báo thành công và ẩn các đối tượng khác
                createFailedObject.SetActive(true);
                foreach (GameObject obj in objectsToDisable)
                {
                    obj.SetActive(false);
                }
            }
            else
            {
                // Xử lý lỗi khi gặp lỗi kết nối hoặc phản hồi không thành công
                Debug.Log("Request failed: " + request.error);
            }
        }
    }
}
