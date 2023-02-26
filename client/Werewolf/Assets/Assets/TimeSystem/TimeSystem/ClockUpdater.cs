using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

public class ClockUpdater : MonoBehaviour
{
	//Referenced Scripts
	public TimeClock clock;

	//Text labels
	public TextMeshProUGUI timeText;

	public TextMeshProUGUI curryearText;
	public TextMeshProUGUI currmonthText;
	public TextMeshProUGUI currdayText;

	//Display Icons
	public GameObject moon;
	public GameObject sun;

	// Start is called before the first frame update
	void Start()
	{
		InvokeRepeating("UpdateClock", 1f, clock.secondSpeed);
	}

	// Update is called once per frame
	void Update()
	{
		
	}
	private void UpdateClock()
	{
		if (clock.night == true)
		{
			sun.SetActive(false);
			moon.SetActive(true);
		}
		else
		{
			sun.SetActive(true);
			moon.SetActive(false);
		}
		curryearText.text = clock.yy.ToString();
		currmonthText.text = clock.month.monthName;
		currdayText.text = clock.days.ToString();		

		string hours = clock.hh.ToString();
		string minutes = clock.mm.ToString();
		if (hours.Length <1 )
		{
			hours = "0"+ clock.hh.ToString();
		}
		if (minutes.Length <= 1)
		{
			minutes = "0" + clock.mm.ToString();
		}
		timeText.text = hours + ":" + minutes;

	}
}
