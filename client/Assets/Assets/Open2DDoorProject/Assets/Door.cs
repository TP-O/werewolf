using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Door : MonoBehaviour
{
    private bool coroutineAllowed, opened;

    private AudioSource audSrc;

    [SerializeField]
    private AudioClip openSound, closeSound;

    private void Start()
    {
        audSrc = GetComponent<AudioSource>();
        coroutineAllowed = true;
        opened = false;
    }

    private void OnMouseDown()
    {
        if (coroutineAllowed)
        {
            StartCoroutine(OpenDoor());
        }
    }

    private IEnumerator OpenDoor()
    {
        coroutineAllowed = false;

        if (!opened)
        {
            audSrc.PlayOneShot(openSound);

            for (float i = 0f; i >= -180f; i -= 10f)
            {
                transform.rotation = Quaternion.Euler(0f, i, 0f);
                yield return new WaitForSeconds(0.01f);
            }
        }
        else if (opened)
        {
            audSrc.PlayOneShot(closeSound);

            for (float i = -180f; i <= 0f; i += 10f)
            {
                transform.rotation = Quaternion.Euler(0f, i, 0f);
                yield return new WaitForSeconds(0.01f);
            }
        }

        coroutineAllowed = true;
        opened = !opened;
    }
}
