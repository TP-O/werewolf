using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.InputSystem;

public class PlayerMovement : MonoBehaviour
{
    [SerializeField] private int speed = 5;
    private Vector2 movement;
    private Rigidbody2D rb;
    private Animator animator;

    public VectorValue startingPosition;

    private void Awake()
    {
        rb = GetComponent<Rigidbody2D>();
        animator = GetComponent<Animator>();
    }

    private void Start()
    {
        transform.position = startingPosition.initialValue;
    }

    private void OnMovement(InputValue value)
    {
        movement = value.Get<Vector2>();
        if (movement.x != 0 || movement.y != 0)
        {
            animator.SetFloat("X", movement.x);
            animator.SetFloat("Y", movement.y);

            animator.SetBool("isWalking", true);
        }
        else
        {
            animator.SetBool("isWalking", false);
        }
            
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
 