using UnityEngine;

public class PlayerController : MonoBehaviour
{
    public string id;
    private Color normalColor = Color.white;
    private Color hiddenColor = Color.black;
    private SpriteRenderer spriteRenderer;

    void Start()
    {
        spriteRenderer = GetComponent<SpriteRenderer>();
        //spriteRenderer.color = normalColor;
    }

    void Update()
    {
        if (id == "")
        {
            spriteRenderer.color = hiddenColor;
        }
        else
        {
            spriteRenderer.color = normalColor;
        }
    }
}
