using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Shooting : MonoBehaviour
{
    public Transform firePoint;
    public GameObject bulletPrefab;
    public float bulletForce = 20f;

    private bool canShoot = true; // Biến để kiểm tra xem có được bắn hay không
    private int shotCount = 0; // Biến đếm số lần bắn

    void Update()
    {
        if (Input.GetButtonDown("Fire1") && canShoot)
        {
            Shoot();
            shotCount++;
            if (shotCount >= 2)
            {
                canShoot = false; // Không cho phép bắn nữa
                shotCount = 0; // Đặt lại biến đếm
                StartCoroutine(WaitForNextShot()); // Bắt đầu coroutine đợi thời gian chờ
            }
        }
    }

    void Shoot()
    {
        Vector2 targetPoint = Camera.main.ScreenToWorldPoint(Input.mousePosition);
        Vector2 firePointPosition = firePoint.position;
        Vector2 direction = (targetPoint - firePointPosition).normalized;

        GameObject bullet = Instantiate(bulletPrefab, firePoint.position, firePoint.rotation);
        Rigidbody2D rb = bullet.GetComponent<Rigidbody2D>();
        rb.AddForce(direction * bulletForce, ForceMode2D.Impulse);
    }

    IEnumerator WaitForNextShot()
    {
        yield return new WaitForSeconds(1f); // Chờ 1 giây
        canShoot = true; // Cho phép bắn tiếp tục
    }
}
