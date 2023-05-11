using System.Collections;
using UnityEngine.UI;
using UnityEngine.EventSystems;
using System.Collections.Generic;
using UnityEngine;
using TMPro;
using System;

public class ChangeInput : MonoBehaviour
{
    public TMP_Text clockText;

    EventSystem system;
    //public Selectable firstInput;
    public Button submitButton;

    void Start()
    {
        system = EventSystem.current;
        //firstInput.Select();
    }

    void Update()
    {
        //clockText.text = DateTime.Now.ToString();

        if (Input.GetKeyDown(KeyCode.Tab) && Input.GetKey(KeyCode.LeftShift))
        {
            Selectable previous = system.currentSelectedGameObject.GetComponent<Selectable>().FindSelectableOnUp();
            if(previous != null)
            {
                previous.Select();
            }    
        }
        else if (Input.GetKeyDown(KeyCode.Tab))
        {
            Selectable next = system.currentSelectedGameObject.GetComponent<Selectable>().FindSelectableOnDown();
            if (next != null)
            {
                next.Select();
            }
        }
        else if (Input.GetKeyDown(KeyCode.Return))
        {
            submitButton.onClick.Invoke();
            Debug.Log("Login successfully!");
        }
    }
}
