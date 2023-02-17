using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.SceneManagement;

public class Door1 : MonoBehaviour
{
    private Player thePlayer;

    public SpriteRenderer theSR;
    public Sprite doorOpenSprite;

    public bool doorOpen, waitingToOpen;

    //public int iLevelToLoad;
    //public string sLevelToLoad;

    //public bool useIntegerToLoadLevel = false;



    // Start is called before the first frame update
    void Start()
    {
        thePlayer = FindObjectOfType<Fox>();
    }

    // Update is called once per frame
    void Update()
    {
        if (waitingToOpen)
        {
            if (Vector3.Distance(thePlayer.followingKey.transform.position, transform.position) < 0.1f)
            {
                waitingToOpen = false;

                doorOpen = true;

                theSR.sprite = doorOpenSprite;

                thePlayer.followingKey.gameObject.SetActive(false);
                thePlayer.followingKey = null;
            }
        }

        //Debug.Log("Qua duoc ma!!" + doorOpen + " - " + Vector3.Distance(thePlayer.transform.position, transform.position) + " - " + Input.GetAxis("Vertical"));

        /*if (doorOpen && Vector2.Distance(thePlayer.transform.position, transform.position) < 1.5f && Input.GetAxis("Vertical") > 0.1f)
        {
            //SceneManager.LoadScene(SceneManager.GetActiveScene().name);

            LoadScene();
        }*/
    }


    private void OnTriggerEnter2D(Collider2D other)
    {
        if (other.tag == "Player")
        {
            if (thePlayer.followingKey != null)
            {
                thePlayer.followingKey.followTarget = transform;
                waitingToOpen = true;

                GameObject collisionGameObject = other.gameObject;
                if (collisionGameObject.name == "Fox")
                {
                    //LoadScene();
                }
            }
        }
    }


    /*void LoadScene()
    {
        if (useIntegerToLoadLevel)
        {
            SceneManager.LoadScene(iLevelToLoad);
        }
        else
        {
            SceneManager.LoadScene(sLevelToLoad);
        }
    }*/
}
