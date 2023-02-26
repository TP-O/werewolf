using System.Collections;
using System.Collections.Generic;
using UnityEngine;


public class TimeClock : MonoBehaviour
{
	public SOMont[] scriptableMonths =  new SOMont[12];
	public SOMont month;

	//TimeClock Properties
	public int yy;
	public int days;
	public int actualMonth;
	public int hh;
	public int mm;
	public bool night = false;	

	public float secondSpeed;

		
	public UnityEngine.Rendering.Universal.Light2D globallight;
	public UnityEngine.Rendering.Universal.Light2D houselight;

	int hoursPassed;
   
	// Start is called before the first frame update
	void Start()
	{
		month = scriptableMonths[1];
		InvokeRepeating("TimePasses", secondSpeed, secondSpeed);
		InvokeRepeating("DaynightSwitch", 2f, 2f);
		
	}
	// Update is called once per frame
	
	void FixedUpdate()
	{
		SwitchLights();
	}
	public void TimePasses() //Sets the IngameTime passing
	{
		mm++;
		if (mm > 59)
		{
			mm = 0;
			hh++;
			hoursPassed++;
			CheckClock();
			
		}
	}
	private void CheckClock() //Checks whether the clock has arrived to time 24 and resets it to00;
	{
		if (hh > 23)
		{
			hh = 0;
			days++;
			if (days > month.dayAmmount)
			{
				days = 1;
				actualMonth++;				
				month = scriptableMonths[actualMonth];				
				if (actualMonth >= 12)
				{
					actualMonth = 0;
					yy++;
				}
			}
		}
	}
	public void SwitchLights() 
	{
		float targetIntensity;
		if (night)
		{
			targetIntensity = 0.1f;
			globallight.intensity = Mathf.Lerp(globallight.intensity, targetIntensity, Time.deltaTime*0.2f);
			houselight.intensity = Mathf.Lerp(houselight.intensity, 1f, Time.deltaTime*1f);
		}
		else
		{
			targetIntensity = 0.95f;
			globallight.intensity = Mathf.Lerp(globallight.intensity, targetIntensity, Time.deltaTime*0.2f);
			houselight.intensity = Mathf.Lerp(houselight.intensity, 0f, Time.deltaTime*1f);
		}
	}

	public void DaynightSwitch() //Switches between day and night
	{
		if (hh < 21 && hh > 5)
		{
			night = false;            
		}
		else
		{
			night = true; 
		}		
	}	
}