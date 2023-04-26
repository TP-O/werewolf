using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class AimTarget : MonoBehaviour
{
    public Rigidbody2D rb;
    public Camera cam;
    public Transform player;
    public Transform gun;

    Vector2 mousePos;

    void Update()
    {
        mousePos = cam.ScreenToWorldPoint(Input.mousePosition);
    }

    void FixedUpdate()
    {
        gun.position = player.position; // Đặt vị trí của súng bằng vị trí của người chơi

        Vector2 lookDir = mousePos - rb.position;
        float angle = Mathf.Atan2(lookDir.y, lookDir.x) * Mathf.Rad2Deg;
        rb.rotation = angle;

        // Kiểm tra xem góc quay của súng có lớn hơn 90 hoặc nhỏ hơn -90 không
        if (angle > 90f || angle < -90f)
        {
            // Nếu có, đổi chiều scale của súng theo trục y
            gun.localScale = new Vector3(1f, -1f, 1f);
        }
        else
        {
            gun.localScale = new Vector3(1f, 1f, 1f);
        }
    }


}