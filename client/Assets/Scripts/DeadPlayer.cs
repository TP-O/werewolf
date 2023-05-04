using UnityEngine;

public class DeadPlayer : MonoBehaviour
{
    private bool isHit = false; // Biến để theo dõi trạng thái của nhân vật

    void Update()
    {
        if (isHit)
        {
            gameObject.SetActive(false); // Nếu bị dính đạn thì ẩn nhân vật đi
        }
    }

    void OnCollisionEnter2D(Collision2D collision)
    {
        if (collision.gameObject.CompareTag("Bullet"))
        {
            isHit = true; // Nếu va chạm với đạn thì gán giá trị true cho biến isHit
        }
    }
}
