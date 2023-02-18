using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.InputSystem;

public class PlayerMovement : MonoBehaviour
{
    [SerializeField] private int speed = 5;
    private Vector2 movement;
    private Rigidbody2D rb;

    public VectorValue startingPosition;

    private void Awake()
    {
        rb = GetComponent<Rigidbody2D>();
    }

    private void Start()
    {
        transform.position = startingPosition.initialValue;
    }

    private void OnMovement(InputValue value)
    {
        movement = value.Get<Vector2>();
    }

    private void FixedUpdate()
    {
        // First way to move character.
        //rb.MovePosition(rb.position + movement * speed * Time.fixedDeltaTime);

        // Second way
        if(movement.x != 0 || movement.y != 0)
        {
            rb.velocity = movement * speed;
        }

        // Third way
        //rb.AddForce(movement * speed);
    }
}
 