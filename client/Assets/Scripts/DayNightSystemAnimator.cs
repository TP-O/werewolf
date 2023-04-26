using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class DayNightSystemAnimator : MonoBehaviour
{
    private Animator animator;
    private string currentState = "Night_idle";
    private bool isNight = true;
    private bool isDay = false;
    private bool isDusk = false;

    public float waitTime = 5f; // thời gian chờ giữa các chuyển đổi state

    private void Awake()
    {
        animator = GetComponent<Animator>();
    }

    private void Start()
    {
        SetState("Night_idle"); // set initial state to Night_idle
        StartCoroutine(TransitionState()); // start coroutine to automatically transition state
    }

    private void SetState(string state)
    {
        switch (state)
        {
            case "Night_idle":
                isNight = true;
                isDay = false;
                isDusk = false;
                break;
            case "Night_Day_transition":
                isNight = false;
                isDay = true;
                isDusk = false;
                break;
            case "Day_Dusk_transition":
                isNight = false;
                isDay = false;
                isDusk = true;
                break;
            case "Dusk_Night_transition":
                isNight = true;
                isDay = false;
                isDusk = false;
                break;
            default:
                Debug.LogWarning($"Invalid state: {state}");
                break;
        }

        animator.SetBool("isDay", isDay);
        animator.SetBool("isDusk", isDusk);
        animator.SetBool("isNight", isNight);

        currentState = state;
    }

    private IEnumerator TransitionState()
    {
        while (true)
        {
            yield return new WaitForSeconds(waitTime); // wait for 20 seconds
            if (currentState == "Night_idle")
            {
                SetState("Night_Day_transition");
            }
            else if (currentState == "Night_Day_transition")
            {
                SetState("Day_Dusk_transition");
            }
            else if (currentState == "Day_Dusk_transition")
            {
                SetState("Dusk_Night_transition");
            }
            else if (currentState == "Dusk_Night_transition")
            {
                SetState("Night_Day_transition");
            }
        }
    }
}
