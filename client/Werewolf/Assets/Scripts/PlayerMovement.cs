using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.InputSystem;

public class PlayerMovement : MonoBehaviour
{
    [SerializeField] private int speed = 5;
    private Vector2 movement;
    private Rigidbody2D rb;

    private void Awake()
    {
        rb = GetComponent<Rigidbody2D>();
    }

    private void OnMovement(InputValue value)
    {
        movement = value.Get<Vector2>();
    }

    private void FixedUpdate()
    {
        // Variant 1
        //rb.MovePosition(rb.position + movement * speed * Time.fixedDeltaTime);

        // Variant 2
        if(movement.x != 0 || movement.y != 0)
        {
            rb.velocity = movement * speed;
        }

        // Variant 3
        //rb.AddForce(movement * speed);
    }
}
