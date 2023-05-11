using UnityEngine;

public class CanvasManager : MonoBehaviour
{
    public static GameObject[] canvas;

    private void Start()
    {
        canvas[0].SetActive(true);
    }
}