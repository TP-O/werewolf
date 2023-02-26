using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public enum DayCycles
{
    Night = 0,
    Day = 1,
    Dusk = 2
}

public class DayNightSystem2D : MonoBehaviour
{
    [Header("Controllers")]
    public UnityEngine.Rendering.Universal.Light2D globalLight;
    //public float cycleCurrentTime = 0;
    public float cycleMaxTime = 5;
    public float colorChangeDuration = 2f;
    public DayCycles dayCycle = DayCycles.Night;

    [Header("Cycle Colors")]
    public Color night;
    public Color day;
    public Color dusk;

    [Header("Objects")]
    public UnityEngine.Rendering.Universal.Light2D[] mapLights;

    void Start()
    {
        dayCycle = DayCycles.Night;
        globalLight.color = night;
        StartCoroutine(ChangeDayCycle());
    }

    IEnumerator ChangeDayCycle()
    {
        while (true)
        {
            yield return new WaitForSeconds(5f);

            switch (dayCycle)
            {
                case DayCycles.Night:
                    StartCoroutine(ChangeColorOverTime(globalLight, night, day, colorChangeDuration));
                    dayCycle = DayCycles.Day;
                    break;

                case DayCycles.Day:
                    StartCoroutine(ChangeColorOverTime(globalLight, day, dusk, colorChangeDuration));
                    dayCycle = DayCycles.Dusk;
                    break;

                case DayCycles.Dusk:
                    StartCoroutine(ChangeColorOverTime(globalLight, dusk, night, colorChangeDuration));
                    dayCycle = DayCycles.Night;
                    break;

                default:
                    break;
            }
        }
    }

    IEnumerator ChangeColorOverTime(UnityEngine.Rendering.Universal.Light2D light, Color startColor, Color endColor, float duration)
    {
        float timeElapsed = 0f;

        while (timeElapsed < duration)
        {
            timeElapsed += Time.deltaTime;
            light.color = Color.Lerp(startColor, endColor, timeElapsed / duration);
            yield return null;
        }

        light.color = endColor;
    }

    void ControlLightMaps(bool status)
    {
        if (mapLights.Length > 0)
            foreach (UnityEngine.Rendering.Universal.Light2D _light in mapLights)
                _light.gameObject.SetActive(status);
    }
}
